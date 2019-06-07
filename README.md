## Firecracker Ignite

Ignite is a Firecracker microVM administration tool, like Docker manages
runC containers.
It builds VM images from OCI images, spin VMs up/down in lightning speed,
and manages multiple VMs efficiently.

The idea is that Ignite makes Firecracker VMs look like Docker containers.
So we can deploy and manage full-blown VM systems just like e.g. Kubernetes workloads.
The images used are Docker images, but instead of running them in a container, the root
filesystem of the image executes as a real VM with a dedicated kernel and `/sbin/init` as
PID 1.

Networking is set up automatically, the VM gets the same IP as any docker
container on the host would.

And Firecracker is **fast**! Building and starting VMs takes just some fraction of a second, or
at most some seconds. With Ignite you can get started with Firecracker in no time!

### Use-cases

With Ignite, Firecracker is now much more accessible for end users, which means the ecosystem
can achieve the next level of momentum due to the easy onboarding path thanks to a docker-like UX.

Although Firecracker was designed with serverless workloads in mind, it can equally well boot a
normal Linux OS, like Ubuntu, Debian or CentOS, running an init system like `systemd`.

Having a super-fast way of spinning up a new VM, with a kernel of choice, running an init system
like `systemd` allows to run system-level applications like the kubelet, which needs to “own” the full system.

This allows for:
* Legacy applications which cannot be containerized (e.g. they need a specific kernel)
  * Alternative, a very new type of application requiring 
* Reproducible, fast testing of system-wide programs (like Weave Net)
* Super fast Kubernetes Cluster Lifecycle with multiple machines (without docker hacks)
* A k8s-managed private VM cloud, on which a layer of k8s container clusters may run

#### Scope

If you want to run _applications_ in **containers** with added _Firecracker isolation_, use
[firecracker-containerd](https://github.com/firecracker-microvm/firecracker-containerd).
Or a similar solution like Kata Containers or gVisor, that are complementary to firecracker-containerd.

Firecracker Ignite, however, is operating at another layer. Ignite isn’t concerned with **containers**
as the primary unit, but whole yet lightweight VMs that integrate with the container landscape.

### How to use

[![asciicast](https://asciinema.org/a/la4p00rcJOBLfW4HPvS8z7riy.svg)](https://asciinema.org/a/la4p00rcJOBLfW4HPvS8z7riy)

Note: At the moment `ignite` needs root privileges on the host to operate,
for certain specific operations (e.g. `mount`). This will change in the future.

```console
$ ignite build weaveworks/ignite-ubuntu:v0.2.0 \
    --name ubuntu-image \
    --import-kernel ubuntu-kernel
$ ignite images
$ ignite kernels
$ ignite run ubuntu-image ubuntu-kernel --name my-vm
$ ignite ps
$ ignite logs my-vm
$ ignite attach my-vm

# Cleanup
$ ignite stop my-vm
$ ignite rm my-vm
$ ignite rmi ubuntu-image
$ ignite rmk ubuntu-kernel
```

When doing `ignite attach`, first press Enter, and then
login with user "root" and password "root".

### CLI documentation

See the [CLI Reference](docs/cli/ignite.md).

### Sample Images

As the upstream `centos:7` and `ubuntu:18.04` images from Docker Hub doesn't
have all the utilities and packages you'd expect in a VM, we have packaged some
reference base images and a sample kernel image to get started quickly.

 - [Kernel Builder Image](images/kernel/Dockerfile)
 - [Ubuntu 18.04 Dockerfile](images/ubuntu/Dockerfile)
 - [CentOS 7 Dockerfile](images/ubuntu/Dockerfile)

### Known limitations

See [REQUIREMENTS.md](REQUIREMENTS.md)

### Maintainers

- Lucas Käldström, @luxas
- Dennis Marttinen, @twelho

### License

Apache 2.0
