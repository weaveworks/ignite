package gitops

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/weaveworks/flux/git"
)

func NewGitDirectory(url, dir, branch string, interval time.Duration) *GitDirectory {
	d := &GitDirectory{
		url: url,
		repo: git.NewRepo(git.Remote{
			URL: url,
		}, git.ReadOnly, git.PollInterval(interval)),
		syncInterval: interval,
		wg:           &sync.WaitGroup{},
		stop:         make(chan struct{}, 1),
		updates:      make(chan GitUpdate),
		dir:          dir,
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
	dir          string
	updates      chan GitUpdate

	branch     string
	lastCommit string
	lastExport string
}

type GitUpdate struct {
	Commit        string
	PersistentDir string
	RealDir       string
}

func (d *GitDirectory) StartLoop() {
	go func() {
		if err := d.repo.Start(d.stop, d.wg); err != nil {
			log.Fatalf("The Git syncing loop terminated with an error %v\n", err)
		}
	}()

	go func() {
		if err := d.checkoutLoop(); err != nil {
			log.Fatalf("The GitOps loop terminated with an error %v\n", err)
		}
	}()
}

func (d *GitDirectory) checkoutLoop() error {
	hasBeenReadyBefore := false
	for {
		time.Sleep(d.syncInterval)

		// check status in order to wait until the repo is ready
		status, err := d.repo.Status()
		if err != nil {
			return err
		}

		if status != git.RepoReady {
			continue
		}

		if !hasBeenReadyBefore {
			log.Printf("Git initialized: A bare clone of repo %q has been made\n", d.url)
			hasBeenReadyBefore = true
		}

		commit, err := d.repo.Revision(context.Background(), d.branch)
		if err != nil {
			return fmt.Errorf("revision error %v", err)
		}

		if commit == d.lastCommit {
			continue
		}

		ex, err := d.repo.Export(context.Background(), d.branch)
		if err != nil {
			return fmt.Errorf("export error %v", err)
		}

		// TODO: Clean up the previous exports after some more commits
		if err := os.RemoveAll(d.dir); err != nil {
			return fmt.Errorf("symlink remove error %v", err)
		}

		if err := os.RemoveAll(d.lastExport); err != nil {
			return fmt.Errorf("lastexport remove error %v", err)
		}

		if err := os.Symlink(ex.Dir(), d.dir); err != nil {
			return fmt.Errorf("symlink error %v", err)
		}

		d.lastCommit = commit
		d.lastExport = ex.Dir()
		d.updates <- GitUpdate{
			Commit:        commit,
			PersistentDir: d.dir,
			RealDir:       ex.Dir(),
		}

		log.Printf("New commit observed on branch %q: %s", d.branch, commit)
	}
}

func (d *GitDirectory) WaitForUpdate() GitUpdate {
	return <-d.updates
}

func (d *GitDirectory) Ready() bool {
	r := d.lastExport != ""
	return r
}
