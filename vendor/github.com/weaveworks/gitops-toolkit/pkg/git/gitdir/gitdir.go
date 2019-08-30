package gitdir

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/gitops-toolkit/pkg/git"
	"github.com/weaveworks/gitops-toolkit/pkg/util"
)

func NewGitDirectory(url, branch string, paths []string, interval time.Duration) *GitDirectory {
	d := &GitDirectory{
		url: url,
		repo: git.NewRepo(git.Remote{
			URL: url,
		}, git.PollInterval(interval), git.Branch(branch)), // git.ReadOnly
		gitConfig: git.Config{
			Branch:    branch,
			Paths:     paths,
			UserName:  "Weave Ignite",
			UserEmail: "support@weave.works",
			SyncTag:   "ignite-gitops",
			NotesRef:  "ignite-gitops",
		},
		syncInterval: interval,
		wg:           &sync.WaitGroup{},
		stop:         make(chan struct{}, 1),
		branch:       branch,
	}
	// Wait for one goroutine, the one syncing the remote
	d.wg.Add(1)
	// Start syncing directly
	d.startLoops()
	return d
}

type GitDirectory struct {
	url          string
	repo         *git.Repo
	wg           *sync.WaitGroup
	stop         chan struct{}
	syncInterval time.Duration

	branch     string
	paths      []string
	gitConfig  git.Config
	checkout   *git.Checkout
	lastCommit string

	lock sync.Mutex
}

func (d *GitDirectory) Dir() string {
	if d.checkout == nil {
		return ""
	}
	return d.checkout.Dir()
}

func (d *GitDirectory) startLoops() {
	go func() {
		log.Debugf("Starting the repo sync...")
		if err := d.repo.Start(d.stop, d.wg); err != nil {
			log.Fatalf("The Git syncing loop terminated with an error %v\n", err)
		}
	}()

	go func() {
		log.Debugf("Starting the checkout loop...")
		if err := d.checkoutLoop(); err != nil {
			log.Fatalf("The GitOps checkout loop terminated with an error %v\n", err)
		}
	}()

	go func() {
		log.Debugf("Starting the commit loop...")
		if err := d.commitLoop(); err != nil {
			log.Fatalf("The GitOps commit loop terminated with an error %v\n", err)
		}
	}()
}

func (d *GitDirectory) checkoutLoop() error {
	// First, loop to clone and do the initial checkout
	log.Info("Initializing the Git repo...")
	for {
		time.Sleep(d.syncInterval / 2)

		// check status in order to wait until the repo is ready
		// tolerate the ClonedOnly and NotCloned errors; just retry
		status, err := d.repo.Status()
		if err != nil && err != git.ErrClonedOnly && err != git.ErrNotCloned {
			return err
		}

		if status != git.RepoReady {
			continue
		}

		log.Infof("Git initialized: A bare clone of repo %q has been made\n", d.url)

		// Clone the checkout for the first time, otherwise just forward the branch
		d.checkout, err = d.repo.Clone(context.Background(), d.gitConfig)
		if err != nil {
			return fmt.Errorf("checkout clone error %v", err)
		}

		commit, err := d.repo.BranchHead(context.Background())
		if err != nil {
			return fmt.Errorf("revision error %v", err)
		}

		// Notify upstream that we now have a commit
		d.observeCommit(commit, true)

		// Now, break out of this loop and go to the next one
		break
	}

	log.Info("Initial clone done, entering the checkout loop...")
	for {
		time.Sleep(d.syncInterval / 2)

		commit, err := d.repo.BranchHead(context.Background())
		if err != nil {
			return fmt.Errorf("revision error %v", err)
		}

		// Just continue looping when the commit hasn't changed
		if commit == d.lastCommit {
			continue
		}

		// Perform the checkout of the new revision
		if err := d.doCheckout(commit); err != nil {
			return err
		}
	}
}

func (d *GitDirectory) doCheckout(commit string) error {
	// Lock the mutex now that we're starting, and unlock it when exiting
	d.lock.Lock()
	defer d.lock.Unlock()

	// If the clone already exists, git fetch the latest contents and checkout the new commit
	if _, err := util.ExecuteCommand("git", "-C", d.Dir(), "fetch"); err != nil {
		return fmt.Errorf("git fetch error %v", err)
	}
	if err := d.checkout.Checkout(context.Background(), commit); err != nil {
		return fmt.Errorf("checkout update error %v", err)
	}

	// Notify upstream that we now have a new commit, and allow writing again
	d.observeCommit(commit, true)
	return nil
}

// observeCommit sets the lastCommit variable so that we know the latest state
func (d *GitDirectory) observeCommit(commit string, userInitiated bool) {
	d.lastCommit = commit
	log.Infof("New commit observed on branch %q: %s. User initiated: %t", d.branch, commit, userInitiated)
}

func (d *GitDirectory) commitLoop() error {
	for {
		time.Sleep(d.syncInterval / 2)

		// Wait for the checkout to exist
		if d.checkout == nil {
			continue
		}

		files, err := d.checkout.ChangedFiles(context.Background(), d.lastCommit)
		if err != nil {
			log.Errorf("couldn't get changed files: %v", err)
			continue
		}

		if len(files) == 0 {
			log.Tracef("no changed files in git repo, nothing to commit...")
			continue
		}

		// Perform the commit
		if err := d.doCommit(); err != nil {
			return err
		}
	}
}

func (d *GitDirectory) doCommit() error {
	// Lock the mutex now that we're starting, and unlock it when exiting
	d.lock.Lock()
	defer d.lock.Unlock()

	// Do a commit and push
	if err := d.checkout.CommitAndPush(context.Background(), git.CommitAction{
		//Author: d.gitConfig.UserName,
		Message: "Update files changed by Ignite",
	}, nil, true); err != nil {
		return fmt.Errorf("git commit and/or push error: %v", err)
	}

	commit, err := d.checkout.HeadRevision(context.Background())
	if err != nil {
		return fmt.Errorf("git revision error: %v", err)
	}

	// Notify upstream that we now have a new commit, and allow writing again
	log.Infof("A new commit with the actual state has been created and pushed to the origin: %q", commit)
	d.observeCommit(commit, false)
	return nil
}

func (d *GitDirectory) Ready() bool {
	r := d.lastCommit != ""
	log.Tracef("git directory ready: %t", r)
	return r
}

func (d *GitDirectory) WaitForClone() {
	for {
		// Check if the Git repo is ready and cloned, otherwise wait
		if d.Ready() {
			break
		}

		time.Sleep(d.syncInterval / 2)
	}
}
