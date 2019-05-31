## Firecracker Ignite

Ignite is a containerized Firecracker microVM administration tool.
It can build VM images, spin VMs up/down and manage multiple VMs efficiently.
Integrates well with cloud-native projects like CNI, containerd and Docker.

### How to use

```console
$ ignite build luxas/ubuntu-base:18.04 \
    --name my-image \
    --import-kernel my-kernel
$ ignite images
$ ignite kernels
$ ignite run my-image my-kernel --name my-vm
$ ignite ps
$ ignite attach my-vm
```
Login with user "root" and password "root".

Note: `sysctl net.bridge.bridge-nf-call-iptables=0` might need to be set on the host for Firecracker.

### Maintainers

- Lucas Käldström, @luxas
- Dennis Marttinen, @twelho

### License

Apache 2.0
