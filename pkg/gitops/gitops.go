package gitops

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/operations/reconcile"
	"github.com/weaveworks/libgitops/pkg/gitdir"
	"github.com/weaveworks/libgitops/pkg/storage/manifest"
)

func RunGitOps(url string, opts gitdir.GitDirectoryOptions) error {
	log.Infof("Starting GitOps loop for repo at %q\n", url)
	log.Info("Whenever changes are pushed to the target branch, Ignite will apply the desired state locally\n")

	// Construct the GitDirectory implementation which backs the storage
	gitDir, err := gitdir.NewGitDirectory(url, opts)
	if err != nil {
		return err
	}
	// TODO: Run gitDir.Cleanup() on SIGINT

	// Wait for the repo to be cloned
	if err := gitDir.WaitForClone(); err != nil {
		return err
	}

	// Construct a manifest storage for the path backed by git
	s, err := manifest.NewTwoWayManifestStorage(gitDir.Dir(), constants.DATA_DIR, scheme.Serializer)
	if err != nil {
		return err
	}

	// TODO: Make the reconcile function signal-aware
	reconcile.ReconcileManifests(s)
	return nil
}
