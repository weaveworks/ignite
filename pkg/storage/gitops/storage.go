package gitops

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage"
)

const dataDir = "/tmp/ignite-gitops"

func NewGitOpsStorage(url, branch string) *GitOpsStorage {
	syncInterval, _ := time.ParseDuration("10s")
	gitRaw := NewGitRawStorage(dataDir, constants.DATA_DIR)
	s := &GitOpsStorage{
		gitRaw:  gitRaw,
		Storage: storage.NewGenericStorage(gitRaw, scheme.Serializer),
		gitDir:  NewGitDirectory(url, dataDir, branch, syncInterval),
		updates: make(chan UpdatedFiles),
	}
	s.gitDir.StartLoop()
	s.startSync()
	return s
}

// GitOpsStorage implements the storage interface for GitOps purposes
type GitOpsStorage struct {
	storage.Storage
	gitRaw  *GitRawStorage
	gitDir  *GitDirectory
	updates chan UpdatedFiles
}

func (s *GitOpsStorage) startSync() {
	go func() {
		for {
			// Whenever the git repo updates, resync the files in the repo
			s.gitDir.WaitForUpdate()
			diff, err := s.gitRaw.Sync()
			if err != nil {
				log.Warnf("An error occured while syncing git state %v\n", err)
				continue
			}
			s.updates <- diff
		}
	}()
}

func (s *GitOpsStorage) WaitForUpdate() UpdatedFiles {
	return <-s.updates
}

func (s *GitOpsStorage) Ready() bool {
	return s.gitDir.Ready()
}
