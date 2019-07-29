package gitdir

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/weaveworks/flux/git"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/util"
)

func NewGitDirectory(url, branch string, paths []string, interval time.Duration) *GitDirectory {
	d := &GitDirectory{
		url: url,
		repo: git.NewRepo(git.Remote{
			URL: url,
		}, git.PollInterval(interval), git.Branch(branch)), // git.ReadOnly
		gitConfig: git.Config{
			Branch: branch,
			Paths: paths,
			UserName: "Weave Ignite",
			UserEmail: "support@weave.works",
			SyncTag:   "ignite-gitops",
			NotesRef:  "ignite-gitops",
		},
		syncInterval: interval,
		wg:           &sync.WaitGroup{},
		stop:         make(chan struct{}, 1),
		updates:      make(chan GitUpdate),
		branch:       branch,
	}
	// Wait for one goroutine, the one syncing the remote
	d.wg.Add(1)
	return d
}

type GitDirectory struct {
	url          string
	repo         *git.Repo
	wg           *sync.WaitGroup
	stop         chan struct{}
	syncInterval time.Duration
	updates      chan GitUpdate

	branch     string
	paths []string
	gitConfig git.Config
	checkout *git.Checkout
	lastCommit string
}

type GitUpdate struct {
	Commit        string
}

func (d *GitDirectory) Dir() string {
	if d.checkout == nil {
		return ""
	}
	return d.checkout.Dir()
}

func (d *GitDirectory) StartLoop() {
	go func() {
		log.Debugf("Starting repo sync...")
		if err := d.repo.Start(d.stop, d.wg); err != nil {
			log.Fatalf("The Git syncing loop terminated with an error %v\n", err)
		}
	}()

	go func() {
		log.Debugf("Starting checkout loop...")
		if err := d.checkoutLoop(); err != nil {
			log.Fatalf("The GitOps loop terminated with an error %v\n", err)
		}
	}()
}

func (d *GitDirectory) checkoutLoop() error {
	hasBeenReadyBefore := false
	for {
		time.Sleep(d.syncInterval / 2)

		// check status in order to wait until the repo is ready
		status, err := d.repo.Status()
		if err != nil && err != git.ErrClonedOnly {
			return err
		}
		
		if status != git.RepoReady {
			continue
		}

		if !hasBeenReadyBefore {
			log.Printf("Git initialized: A bare clone of repo %q has been made\n", d.url)
			hasBeenReadyBefore = true
		}

		commit, err := d.repo.BranchHead(context.Background())
		if err != nil {
			return fmt.Errorf("revision error %v", err)
		}

		if commit == d.lastCommit {
			continue
		}

		if d.checkout == nil {
			// Clone the checkout for the first time, otherwise just forward the branch
			d.checkout, err = d.repo.Clone(context.Background(), d.gitConfig)
			if err != nil {
				return fmt.Errorf("checkout clone error %v", err)
			}
			
		} else {
			// If the clone already exists, git fetch the latest contents and checkout the new commit
			if _, err := util.ExecuteCommand("git", "-C", d.Dir(), "fetch"); err != nil {
				return fmt.Errorf("git fetch error %v", err)
			}
			if err := d.checkout.Checkout(context.Background(), commit); err != nil {
				return fmt.Errorf("checkout update error %v", err)
			}
		}

		d.lastCommit = commit
		d.updates <- GitUpdate{
			Commit:        commit,
		}

		log.Printf("New commit observed on branch %q: %s", d.branch, commit)
	}
}

func (d *GitDirectory) WaitForUpdate() GitUpdate {
	return <-d.updates
}

func (d *GitDirectory) Ready() bool {
	r := d.checkout != nil
	log.Debugf("git directory ready: %t", r)
	return r
}
