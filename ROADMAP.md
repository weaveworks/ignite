# Roadmap

This is a provisional roadmap for what we'd like to do in the coming releases.

 - CNI support, when using Docker (similar to what the kubelet does)
 - Prometheus metrics exposed over a socket (converted from the FIFO Firecracker supplies)
 - containerd support, do not hard-code this to docker only. This will add Kubernetes support.
 - Parallelized internal architecture for way better performance
 - Factor out the `ignite container` command to its own binary for less overhead when running VMs (`ignite-spawn`)
 - Take full advantage of Devicemapper's thin-provisioning feature
 - Re-architect Ignite, and split it up into four parts:
    - `ignite`, keeping the same CLI as we have now
    - `ignite-containerd`, a containerd plugin for running Ignite VMs
    - `ignite-snapshotter`, a containerd snapshotter plugin for building VM disks out of OCI images
    - `ignite-spawn`, a small wrapper binary around Firecracker handling networking and DHCP translation
