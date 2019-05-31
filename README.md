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

### How to download

Due to this being a private repo, you need to download the binaries manually by going to the releases
page and clicking the `ignite` binary, and `ignite.tar`. `ignite.tar` is the Docker image Ignite uses
for executing the Firecracker process in a container.

When you have downloaded the two files, do the following:

```
cd Downloads
mv ignite /usr/local/bin
docker load -i ignite.tar
```

### Build from source

The only build requirement is Docker.

```
make binary
make image
```

### Maintainers

- Lucas Käldström, @luxas
- Dennis Marttinen, @twelho

### License

Apache 2.0
