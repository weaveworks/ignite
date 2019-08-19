package gitops

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/gitops/gitdir"
	"github.com/weaveworks/ignite/pkg/operations/reconcile"
	"github.com/weaveworks/ignite/pkg/storage/manifest"
)

const syncInterval = 10 * time.Second

func RunGitOps(url, branch string, paths []string) error {
	log.Infof("Starting GitOps loop for repo at %q\n", url)
	log.Infof("Whenever changes are pushed to the %s branch, Ignite will apply the desired state locally\n", branch)

	// Construct the GitDirectory implementation which backs the storage
	gitDir := gitdir.NewGitDirectory(url, branch, paths, syncInterval)

	// Wait for the repo to be cloned
	gitDir.WaitForClone()

	// Construct a manifest storage for the path backed by git
	s, err := manifest.NewTwoWayManifestStorage(gitDir.Dir(), constants.DATA_DIR, scheme.Serializer)
	if err != nil {
		return err
	}

	// TODO: Make the reconcile function signal-aware
	reconcile.ReconcileManifests(s)
	return nil
}
