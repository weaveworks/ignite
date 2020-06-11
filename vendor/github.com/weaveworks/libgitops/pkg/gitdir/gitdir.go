package gitdir

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/fluxcd/toolkit/pkg/ssh/knownhosts"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	userName  = "Weave libgitops"
	userEmail = "support@weave.works"

	defaultBranch   = "master"
	defaultRemote   = "origin"
	defaultInterval = 30 * time.Second
	defaultTimeout  = 1 * time.Minute
)

type GitDirectoryOptions struct {
	// Options
	Branch   string        // default "master"
	Interval time.Duration // default 30s
	Timeout  time.Duration // default 1m
	// TODO: Support folder prefixes

	// Authentication
	// For HTTPS basic auth. The password should be e.g. a GitHub access token
	Username, Password *string
	// For Git SSH protocol. This is the bytes of e.g. ~/.ssh/id_rsa, given that ~/.ssh/id_rsa.pub is
	// registered with and trusted by the Git provider.
	IdentityFileContent []byte

	// The file content (in bytes) of the known_hosts file to use for remote (e.g. GitHub) public key verification
	// If you want to use the default git CLI behavior, populate this byte slice with contents from
	// ioutil.ReadFile("~/.ssh/known_hosts").
	KnownHostsFileContent []byte
}

func (o *GitDirectoryOptions) Default() {
	if o.Branch == "" {
		o.Branch = defaultBranch
	}
	if o.Interval == 0 {
		o.Interval = defaultInterval
	}
	if o.Timeout == 0 {
		o.Timeout = defaultTimeout
	}
}

// Create a new GitDirectory implementation
func NewGitDirectory(url string, opts GitDirectoryOptions) (*GitDirectory, error) {
	log.Info("Initializing the Git repo...")

	// Default the options
	opts.Default()

	// Create a temporary directory for the clone
	tmpDir, err := ioutil.TempDir("", "libgitops")
	if err != nil {
		return nil, err
	}
	log.Debugf("Created temporary directory for the git clone at %q", tmpDir)

	// Create the struct
	d := &GitDirectory{
		url:      url,
		branch:   opts.Branch,
		interval: opts.Interval,
		timeout:  opts.Timeout,

		cloneDir:  tmpDir,
		auth:      nil,
		readwrite: false, // only switch to read-write if we've got the right credentials
	}
	// Set up the parent context for this class. d.cancel() is called only at Cleanup()
	d.ctx, d.cancel = context.WithCancel(context.Background())

	// Parse the endpoint URL
	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, err
	}

	// Choose authentication method based on the credentials
	switch ep.Protocol {
	case "ssh":
		// If we haven't got the right credentials, just continue in read-only mode
		if len(opts.IdentityFileContent) == 0 || len(opts.KnownHostsFileContent) == 0 {
			break
		}
		pk, err := ssh.NewPublicKeys("git", opts.IdentityFileContent, "")
		if err != nil {
			return nil, err
		}
		callback, err := knownhosts.New(opts.KnownHostsFileContent)
		if err != nil {
			return nil, err
		}
		pk.HostKeyCallback = callback
		d.auth = pk
		d.readwrite = true
	case "https":
		// If we haven't got the right credentials, just continue in read-only mode
		if opts.Username == nil || opts.Password == nil {
			break
		}
		d.auth = &http.BasicAuth{
			Username: *opts.Username,
			Password: *opts.Password,
		}
		d.readwrite = true
	case "file":
		d.readwrite = true // assuming enough file privileges to access the repo
	default:
		return nil, fmt.Errorf("unsupported endpoint scheme: %s for URL %q", ep.Protocol, url)
	}
	log.Trace("URL endpoint parsed and authentication method chosen")

	if d.readwrite {
		log.Infof("Running in read-write mode, will commit back current status to the repo")
	} else {
		log.Infof("Running in read-only mode, won't write status back to the repo")
	}

	// Start (non-blocking) syncing goroutines directly
	d.startLoops()
	return d, nil
}

// GitDirectory is an implementation which keeps a directory
type GitDirectory struct {
	// user-specified options
	url      string
	branch   string
	interval time.Duration
	timeout  time.Duration
	auth     transport.AuthMethod

	// the temporary directory used for the clone
	cloneDir string
	// whether we're operating in read-write or read-only mode
	readwrite bool

	// go-git objects. wt is the worktree of the repo, persistent during the lifetime of repo.
	repo *git.Repository
	wt   *git.Worktree

	// latest known commit to the system
	lastCommit string

	// the context and its cancel function for the lifetime of this struct (until Cleanup())
	ctx    context.Context
	cancel context.CancelFunc
	// the lock for git operations (so pushing and pulling aren't done simultaneously)
	lock sync.Mutex
}

func (d *GitDirectory) Dir() string {
	return d.cloneDir
}

func (d *GitDirectory) startLoops() {
	go d.checkoutLoop()
	if d.readwrite {
		go d.commitLoop()
	}
}

func (d *GitDirectory) checkoutLoop() {
	log.Info("Starting the checkout loop...")

	// First, clone the repo
	if err := d.clone(); err != nil {
		log.Fatalf("Failed to clone git repo: %v", err)
	}

	wait.NonSlidingUntilWithContext(d.ctx, func(_ context.Context) {

		log.Trace("checkoutLoop: Will perform pull operation")
		// Perform a pull & checkout of the new revision
		if err := d.doPull(); err != nil {
			log.Errorf("checkoutLoop: git pull failed with error: %v", err)
			return
		}

	}, d.interval)
	log.Info("Exiting the checkout loop...")
}

func (d *GitDirectory) commitLoop() {
	log.Info("Starting the commit loop...")
	wait.NonSlidingUntilWithContext(d.ctx, func(_ context.Context) {

		// Wait for the checkout to exist
		if d.wt == nil {
			log.Tracef("commitLoop: Waiting for the clone to exist")
			return
		}

		log.Trace("commitLoop: Will perform commit operation, if any")

		// Perform the commit
		if err := d.doCommit(); err != nil {
			log.Errorf("checkoutLoop: git commit & push failed with error:", err)
			return
		}

	}, d.interval)
	log.Info("Exiting the commit loop...")
}

func (d *GitDirectory) clone() error {
	// Lock the mutex now that we're starting, and unlock it when exiting
	d.lock.Lock()
	defer d.lock.Unlock()

	log.Infof("Starting to clone the repository %s with timeout %s", d.url, d.timeout)
	// Do a clone operation to the temporary directory, with a timeout
	err := d.withTimeout(func(ctx context.Context) error {
		var err error
		d.repo, err = git.PlainCloneContext(ctx, d.Dir(), false, &git.CloneOptions{
			URL:           d.url,
			Auth:          d.auth,
			RemoteName:    defaultRemote,
			ReferenceName: plumbing.NewBranchReferenceName(d.branch),
			SingleBranch:  true,
			NoCheckout:    false,
			//Depth:             1, // ref: https://github.com/src-d/go-git/issues/1143
			RecurseSubmodules: 0,
			Progress:          nil,
			Tags:              git.NoTags,
		})
		return err
	})
	// Handle errors
	switch err {
	case nil:
		// no-op, just continue.
	case context.DeadlineExceeded:
		return fmt.Errorf("git clone operation took longer than deadline %s", d.timeout)
	case context.Canceled:
		log.Tracef("context was cancelled")
		return nil // if Cleanup() was called, just exit the goroutine
	default:
		return fmt.Errorf("git clone error: %v", err)
	}

	// Populate the worktree pointer
	d.wt, err = d.repo.Worktree()
	if err != nil {
		return fmt.Errorf("git get worktree error: %v", err)
	}

	// Get the latest HEAD commit and report it to the user
	ref, err := d.repo.Head()
	if err != nil {
		return err
	}

	d.observeCommit(ref.Hash(), true)
	return nil
}

func (d *GitDirectory) doPull() error {
	// Lock the mutex now that we're starting, and unlock it when exiting
	d.lock.Lock()
	defer d.lock.Unlock()

	// Perform the git pull operation using the timeout
	err := d.withTimeout(func(ctx context.Context) error {
		log.Trace("checkoutLoop: Starting pull operation")
		return d.wt.PullContext(ctx, &git.PullOptions{
			Auth:         d.auth,
			SingleBranch: true,
		})
	})
	// Handle errors
	switch err {
	case nil, git.NoErrAlreadyUpToDate:
		// no-op, just continue. Allow the git.NoErrAlreadyUpToDate error
	case context.DeadlineExceeded:
		return fmt.Errorf("git pull operation took longer than deadline %s", d.timeout)
	case context.Canceled:
		log.Tracef("context was cancelled")
		return nil // if Cleanup() was called, just exit the goroutine
	default:
		return fmt.Errorf("failed to pull: %v", err)
	}

	log.Trace("checkoutLoop: Pulled successfully")

	// get current head
	ref, err := d.repo.Head()
	if err != nil {
		return err
	}

	// check if we changed commits
	if d.lastCommit != ref.Hash().String() {
		// Notify upstream that we now have a new commit, and allow writing again
		d.observeCommit(ref.Hash(), true)
	}

	return nil
}

// observeCommit sets the lastCommit variable so that we know the latest state
func (d *GitDirectory) observeCommit(commit plumbing.Hash, userInitiated bool) {
	d.lastCommit = commit.String()
	log.Infof("New commit observed on branch %q: %s. User initiated: %t", d.branch, commit, userInitiated)
}

func (d *GitDirectory) doCommit() error {
	// Lock the mutex now that we're starting, and unlock it when exiting
	d.lock.Lock()
	defer d.lock.Unlock()

	s, err := d.wt.Status()
	if err != nil {
		return fmt.Errorf("couldn't get status: %v", err)
	}
	if s.IsClean() {
		log.Debugf("No changed files in git repo, nothing to commit...")
		return nil
	}

	// Do a commit and push
	log.Debug("commitLoop: Committing all local changes")
	hash, err := d.wt.Commit("Update files changed by libgitops", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  userName,
			Email: userEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("git commit and/or push error: %v", err)
	}

	// Perform the git push operation using the timeout
	err = d.withTimeout(func(ctx context.Context) error {
		log.Debug("commitLoop: Will push with timeout")
		return d.repo.PushContext(ctx, &git.PushOptions{
			Auth: d.auth,
		})
	})
	// Handle errors
	switch err {
	case nil, git.NoErrAlreadyUpToDate:
		// no-op, just continue. Allow the git.NoErrAlreadyUpToDate error
	case context.DeadlineExceeded:
		return fmt.Errorf("git push operation took longer than deadline %s", d.timeout)
	case context.Canceled:
		log.Tracef("context was cancelled")
		return nil // if Cleanup() was called, just exit the goroutine
	default:
		return fmt.Errorf("failed to push: %v", err)
	}

	// Notify upstream that we now have a new commit, and allow writing again
	log.Infof("A new commit with the actual state has been created and pushed to the origin: %q", hash)
	d.observeCommit(hash, false)
	return nil
}

func (d *GitDirectory) withTimeout(fn func(context.Context) error) error {
	// Create a new context with a timeout. The push operation either succeeds in time, times out,
	// or is cancelled by Cleanup(). In case of a successful run, the context is always cancelled afterwards.
	ctx, cancel := context.WithTimeout(d.ctx, d.timeout)
	defer cancel()

	// Run the function using the context and cancel directly afterwards
	fnErr := fn(ctx)

	// Return the context error, if any, first so deadline/cancel signals can propagate.
	// Otherwise passthrough the error returned from the function.
	if ctx.Err() != nil {
		log.Debugf("operation context yielded error %v to be returned. Function error was: %v", ctx.Err(), fnErr)
		return ctx.Err()
	}
	return fnErr
}

// Ready signals whether the git clone is ready to use
func (d *GitDirectory) Ready() bool {
	r := d.lastCommit != ""
	log.Tracef("git directory ready: %t", r)
	return r
}

// WaitForClone waits until Ready() is true. In case Cleanup() is called before the repo was cloned,
// WaitForClone() returns an error.
func (d *GitDirectory) WaitForClone() error {
	// Create a new context for this wait operation, that is cancelled either by the git clone
	// becoming ready, or Cleanup() is called.
	ctx, cancel := context.WithCancel(d.ctx)

	// Use this flag to determine if the context was cancelled by the parent, or if we got the ready signal
	becameReady := false

	// Start the wait loop using the context we created
	log.Trace("WaitForClone: Starting wait loop")
	wait.NonSlidingUntilWithContext(ctx, func(_ context.Context) {

		// Check if the Git repo is ready and cloned, otherwise wait
		if d.Ready() {
			// set the ready flag to true, cancel this operation & return
			log.Trace("WaitForClone: Got ready signal")
			becameReady = true
			cancel()
			return
		}
	}, 3*time.Second)

	// Signal in the return error what the outcome was. A nil error means the clone is ready to use
	if !becameReady {
		return fmt.Errorf("git clone didn't complete, operation was cancelled before completion")
	}
	return nil
}

// Cleanup cancels running goroutines and operations, and removes the temporary clone directory
func (d *GitDirectory) Cleanup() error {
	// Cancel the context for the two running goroutines, and any possible long-running operations
	d.cancel()

	// Remove the temporary directory
	if err := os.RemoveAll(d.Dir()); err != nil {
		log.Errorf("Failed to clean up temp git directory: %v", err)
		return err
	}
	return nil
}
