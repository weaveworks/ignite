package util

import (
	"os/exec"

	"github.com/weaveworks/ignite/pkg/runtime"
)

// RmiDocker removes an image from docker content store.
func RmiDocker(img string) {
	_, _ = exec.Command(
		"docker",
		"rmi", img,
	).CombinedOutput()
}

// RmiContainerd removes an image from containerd content store.
func RmiContainerd(img string) {
	_, _ = exec.Command(
		"ctr", "-n", "firecracker",
		"image", "rm", img,
	).CombinedOutput()
}

// rmiCompletely removes a given image completely, from ignite image store and
// runtime image store.
func RmiCompletely(img string, cmd *Command, rt runtime.Name) {
	// Remote from ignite content store.
	_, _ = cmd.New().
		With("image", "rm", img).
		Cmd.CombinedOutput()

	// Remove from runtime content store.
	switch rt {
	case runtime.RuntimeContainerd:
		RmiContainerd(img)
	case runtime.RuntimeDocker:
		RmiDocker(img)
	}
}
