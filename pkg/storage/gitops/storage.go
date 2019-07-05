package gitops

import (
	"log"
	"time"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage"
)

const dataDir = "/tmp/ignite-gitops"

func NewGitOpsStorage(url, branch string) *GitOpsStorage {
	syncInterval, _ := time.ParseDuration("10s")
	raw := NewGitRawStorage(dataDir, constants.DATA_DIR)
	s := &GitOpsStorage{
		raw:     raw,
		Storage: storage.NewGenericStorage(raw, scheme.Serializer),
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
	raw     *GitRawStorage
	gitDir  *GitDirectory
	updates chan UpdatedFiles
}

func (s *GitOpsStorage) startSync() {
	go func() {
		for {
			// Whenever the git repo updates, resync the files in the repo
			s.gitDir.WaitForUpdate()
			diff, err := s.raw.Sync()
			if err != nil {
				// TODO: Make a real warning
				log.Printf("[WARNING] An error occured while syncing git state %v\n", err)
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

/*
// Get populates the pointer to the Object given, based on the file content
func (s *GitOpsStorage) Get(obj meta.Object) error {
	return s.gs.Get(obj)
}

// GetByID returns a new Object for the resource at the specified kind/uid path, based on the file content
func (s *GitOpsStorage) GetByID(kind string, uid meta.UID) (meta.Object, error) {
	return s.gs.GetByID(kind, uid)
}

// Set saves the Object to disk. If the object does not exist, the
// ObjectMeta.Created field is set automatically
func (s *GitOpsStorage) Set(obj meta.Object) error {
	return fmt.Errorf("not implemented")
}

// Delete removes an object from the storage
func (s *GitOpsStorage) Delete(kind string, uid meta.UID) error {
	return fmt.Errorf("not implemented")
}

// List lists objects for the specific kind
func (s *GitOpsStorage) List(kind string) ([]meta.Object, error) {
	return s.gs.List(kind)
}

// ListMeta lists all objects' APIType representation. In other words,
// only metadata about each object is unmarshalled (uid/name/kind/apiVersion).
// This allows for faster runs (no need to unmarshal "the world"), and less
// resource usage, when only metadata is unmarshalled into memory
func (s *GitOpsStorage) ListMeta(kind string) (meta.APITypeList, error) {
	return s.gs.ListMeta(kind)
}

// GetCache gets a new Cache implementation for the specified kind
func (s *GitOpsStorage) GetCache(kind string) (storage.Cache, error) {
	return s.gs.GetCache(kind)
}*/
