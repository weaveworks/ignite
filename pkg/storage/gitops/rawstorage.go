package gitops

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/util"

	"sigs.k8s.io/yaml"
)

var splitDirsRegex = regexp.MustCompile(`/?([a-z0-9]+)(/[a-z0-9]*)*`)

func NewGitRawStorage(gitDir, underlyingDir string) *GitRawStorage {
	return &GitRawStorage{
		gitDir: gitDir,
		gitPathPrefixes: map[string]bool{ // we only check in VM state into git atm
			storage.KeyForKind(api.KindVM): true,
		},
		passthrough: storage.NewDefaultRawStorage(underlyingDir),
	}
}

type UpdateType string

const (
	UpdateTypeCreated UpdateType = "Created"
	UpdateTypeChanged UpdateType = "Changed"
	UpdateTypeDeleted UpdateType = "Deleted"
)

type UpdatedFiles []*UpdatedFile

type UpdatedFile struct {
	GitPath  string
	Type     UpdateType
	Checksum string
	APIType  *meta.APIType
}

type GitRawStorage struct {
	// directory that is managed by git
	gitDir string
	// keyFileMap maps the virtual key path to real file paths in the repo
	keyFileMap map[string]*UpdatedFile
	// byKind maps a kind to many virtual key paths for the storage impl
	byKind map[string][]string
	// gitPathPrefixes define the path prefixes we want to store in the git repo
	gitPathPrefixes map[string]bool
	// passthrough defines the underlying storage for Kinds we don't care about in git
	passthrough storage.RawStorage
}

func (r *GitRawStorage) Sync() (UpdatedFiles, error) {
	// provide empty placeholders for new data, overwrite .keyFileMap and .byKind in the end
	newKeyFileMap := map[string]*UpdatedFile{}
	newByKind := map[string][]string{}
	// a slice of files that
	diff := UpdatedFiles{}
	// walk the git repo
	dirToWalk := r.gitDir
	if !strings.HasSuffix(dirToWalk, "/") {
		// filepath.Walk needs a trailing slash to start traversing the directory
		dirToWalk += "/"
	}

	err := filepath.Walk(dirToWalk, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Don't traverse into the .git folder
			if info.Name() == ".git" {
				return filepath.SkipDir
			}

			// continue traversing
			return nil
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".json") {
			obj := meta.NewAPIType()
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// The yaml package supports both YAML and JSON
			if err := yaml.Unmarshal(content, obj); err != nil {
				return err
			}

			// Ignore unknown API objects to Ignite (e.g. Kubernetes manifests)
			if !scheme.Scheme.Recognizes(obj.GroupVersionKind()) {
				log.Debugf("Ignoring file with API version %s and kind %s", obj.APIVersion, obj.Kind)
				return nil
			}

			// Require the UID field to be set
			if len(obj.GetUID()) == 0 {
				log.Infof("Ignoring %s at path %q that does not have .metadata.uid set.", obj.GetKind(), r.gitRelativePath(path))
				return nil
			}

			// Require the Name field to be set
			if len(obj.GetName()) == 0 {
				log.Infof("Ignoring %s at path %q that does not have .metadata.name set.", obj.GetKind(), r.gitRelativePath(path))
				return nil
			}

			keyPath := storage.KeyForUID(obj.GetKind(), obj.GetUID())
			kindKey := storage.KeyForKind(obj.GetKind())

			f := &UpdatedFile{
				GitPath:  path,
				Checksum: sha256sum(content),
				APIType:  obj,
			}
			newKeyFileMap[keyPath] = f
			newByKind[kindKey] = append(newByKind[kindKey], keyPath)
			log.Debugf("Stored file info %v at path %q and parent kind %q", *f, keyPath, kindKey)

			// calculate a diff if files change
			if prevFile, ok := r.keyFileMap[keyPath]; ok {
				// file existed already on the last check
				if prevFile.Checksum != f.Checksum {
					f.Type = UpdateTypeChanged
					diff = append(diff, f)
				}
			} else {
				// file did not exist before, hence an addition
				f.Type = UpdateTypeCreated
				diff = append(diff, f)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// range through all the previous files, and detect deletions
	for keyPath, oldFile := range r.keyFileMap {
		if _, ok := newKeyFileMap[keyPath]; !ok {
			// the file existed in the last iteration, but not now; a delete
			// as oldFile will be removed from memory in the map overwrites below
			// we need to create a net-new object, with a copy of the previous APIType
			newAPIType := oldFile.APIType.DeepCopy()
			diff = append(diff, &UpdatedFile{
				Type:     UpdateTypeDeleted,
				GitPath:  oldFile.GitPath,
				Checksum: oldFile.Checksum,
				APIType:  newAPIType,
			})
		}
	}

	r.keyFileMap = newKeyFileMap
	r.byKind = newByKind

	for _, file := range diff {
		gitFilePath := r.gitRelativePath(file.GitPath)
		action := strings.ToLower(string(file.Type))
		log.Printf("File %q was %s. It describes a %s with UID %q and name %q\n", gitFilePath, action, file.APIType.Kind, file.APIType.GetUID(), file.APIType.GetName())
	}

	if len(diff) == 0 {
		log.Println("No changes in the relevant YAML or JSON files")
	}

	return diff, nil
}

func sha256sum(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (r *GitRawStorage) gitRelativePath(fullPath string) string {
	return strings.TrimPrefix(fullPath, r.gitDir+"/")
}

func (r *GitRawStorage) realPath(key string) string {
	// The "/" prefix is enforced
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	info, ok := r.keyFileMap[key]
	if !ok {
		log.Debugf("GitRawStorage.realPath returned an empty string for key %s", key)
		return ""
	}

	return info.GitPath
}

func (r *GitRawStorage) shouldPassthrough(key string) bool {
	// firstDirName is e.g. "vm" when key is "/vm/foobar123"
	firstDirName := splitDirsRegex.FindStringSubmatch(key)[1]
	// check if this kind should be managed by git. if it's git-managed return false
	_, ok := r.gitPathPrefixes[storage.KeyForKind(meta.Kind(firstDirName))]
	return !ok
}

func (r *GitRawStorage) Read(key string) ([]byte, error) {
	log.Debugf("GitRawStorage.Read: %q", key)
	if r.shouldPassthrough(key) {
		return r.passthrough.Read(key)
	}

	file := r.realPath(key)
	return ioutil.ReadFile(file)
}

func (r *GitRawStorage) Exists(key string) bool {
	log.Debugf("GitRawStorage.Exists: %q", key)
	if r.shouldPassthrough(key) {
		return r.passthrough.Exists(key)
	}

	file := r.realPath(key)
	return util.FileExists(file)
}

func (r *GitRawStorage) Write(key string, content []byte) error {
	log.Debugf("GitRawStorage.Write: %q", key)
	// Write always writes to the underlying (expected) place, and to Git
	if err := r.passthrough.Write(key, content); err != nil {
		return err
	}

	// If this should not be stored in Git, return at this point
	if r.shouldPassthrough(key) {
		return nil
	}

	// Do a normal write to the git-backed file.
	file := r.realPath(key)
	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		return err
	}

	// TODO: Do a git commit here!
	return nil
}

func (r *GitRawStorage) Delete(key string) error {
	log.Debugf("GitRawStorage.Delete: %q", key)
	// Delete always deletes in the underlying (expected) place, and in Git
	if err := r.passthrough.Delete(key); err != nil {
		return err
	}

	// also delete the git-backed file.
	file := r.realPath(key)
	if len(file) == 0 {
		// the source of the delete seem to have been the git repo itself
		// this happens when someone deletes a file from git, then the loop
		// notices resource X should be removed, and issues a Delete() request
		// at the storage. At this point the file does not exist in git anymore,
		// so it's safe to just exit quickly here
		return nil
	}

	return os.RemoveAll(file)
	// TODO: Do a git commit here!
}

func (r *GitRawStorage) List(parentKey string) ([]string, error) {
	log.Debugf("GitRawStorage.List: %q", parentKey)
	if r.shouldPassthrough(parentKey) {
		return r.passthrough.List(parentKey)
	}

	return r.byKind[parentKey], nil
}
