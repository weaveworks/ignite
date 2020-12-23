package source

import (
	"os/exec"

	containerderr "github.com/containerd/containerd/errdefs"
)

// TarExtract extracts all files from a source to a directory
func TarExtract(src Source, dir string, args ...string) error {
	args = append([]string{"-x", "-C", dir}, args...)
	tarCmd := exec.Command("tar", args...)
	reader, err := src.Reader()
	if err != nil {
		return err
	}
	defer reader.Close()

	tarCmd.Stdin = reader
	if err := tarCmd.Start(); err != nil {
		return err
	}

	if err := tarCmd.Wait(); err != nil {
		return err
	}

	if err = src.Cleanup(); err != nil {
		// Ignore the cleanup error if the resource no longer exists.
		if !containerderr.IsNotFound(err) {
			return err
		}
	}
	return nil
}
