# Roadmap

This is a provisional roadmap for what we'd like to do in the coming releases.

- Split the `ignite` CLI into a client-server model (using e.g. `gRPC`)
  - The CLI should only be a thin wrapper that talks to `ignited`
  - With this `ignite` can be run without root
  - `ignited` will be run with `root` privileges, or in a container with capabilities specifically set
- Provide deb/rpm packages for an easier installation
- Add Virtual Kubelet support to `ignited`
  - `ignited` will register as a Virtual Kubelet in the target Kubernetes cluster
  - The `VM` API type will be register as a `CustomResourceDefinition`
- Use device-mapper Thin Provisioning for layering image -> kernel -> resize -> writable overlay
  - We might be able to utilize/vendor in [containerd's devicemapper snapshotter](https://godoc.org/github.com/containerd/containerd/snapshots/devmapper)
- Define what's in and out of scope for Ignite clearly, e.g.
  - Supporting to restart VMs or not
  - Supporting multiple network interfaces or not
  - Create one architecture diagram more and a design document
- Parallelized internal architecture for better performance
- Generate OpenAPI documentation and specifications
- Add support for CSI volumes
