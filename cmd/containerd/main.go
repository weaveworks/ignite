package main

import (
	"github.com/weaveworks/ignite/pkg/containerd"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	containerd.Main()
}
