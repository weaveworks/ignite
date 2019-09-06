package containerd

import (
	_ "github.com/containerd/containerd/metrics/cgroups"
	_ "github.com/containerd/containerd/runtime/v1/linux"
	_ "github.com/containerd/containerd/runtime/v2"
	_ "github.com/containerd/containerd/runtime/v2/runc/options"
	_ "github.com/containerd/containerd/snapshots/native"
	_ "github.com/containerd/containerd/snapshots/overlay"
)
