# Roadmap

This is a provisional roadmap for what we'd like to do in the coming releases.

 - containerd support, do not hard-code this to docker only. This will add Kubernetes support.
 - Parallelized internal architecture for better performance
 - Patch support to `pkg/storage` to allow for race-condition-safe writes
 - Write support for the GitOps mode, so it can write status back to the repo
 - Add automated CI testing
 - Provide deb/rpm packages
 - Generate more/better API documentation
 - Add support for CSI volumes
 - Integrate Ignite with [Footloose](https://github.com/weaveworks/footloose)
 - Take full advantage of Devicemapper's thin-provisioning feature
 - Re-architect Ignite, and split it up into four parts:
    - `ignite`, keeping the same CLI as we have now
    - `ignite-containerd`, a containerd plugin for running Ignite VMs
    - `ignite-snapshotter`, a containerd snapshotter plugin for building VM disks out of OCI images
    - `ignite-spawn`, a small wrapper binary around Firecracker handling networking and DHCP translation
