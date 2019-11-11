# FAQ

A collection of Frequently Asked Questions about Ignite:

## Q: Can I use Ignite as CRI runtime for kubelet?

No, you can't. Ignite isn't designed for running containers, hence it cannot work as a CRI runtime.

Ignite runs VMs instead. In the future, we envision Ignite to (maybe) be able to run VMs (not containers)
based off Kubernetes Pods using some special annotations. This would however (most likely) be done as a containerd
plugin (lower in the stack than CRI).

## Q: What is the difference between {Kata Containers, gVisor, fc-containerd} and Ignite?

Kata Containers, gVisor, and firecracker-containerd run _containers_, and Ignite runs _VMs_.

Kata can integrate with Firecracker, but the value add there is more isolation, as the container is
spawned inside of a minimal Firecracker VM.

[firecracker-containerd](https://github.com/firecracker-microvm/firecracker-containerd) enables you to
do the same as Kata; add isolation for a container; but this time in a bit more lightweight manner, as a
containerd plugin.

gVisor acts as a gatekeeper between your application in a container, and the kernel. gVisor emulates the
kernel syscalls, and based on if they are "safe" or not, passes them though to the underlying kernel, or
performs a similar operation. gVisor's value-add is the same as the above, added isolation for containers.

Ignite however, uses the rootfs from an OCI image, and runs that content as a real VM. Inside of the
Firecracker VM spawned, there are no extra containers running (unless the user installs a container
runtime).

## Q: Why does Ignite require KVM?

Firecracker is a KVM implementation, and uses KVM to manage and virtualize the VM.

## Q: Why does Ignite require to be root?

In order to prepare the VM filesystem, Ignite needs to create a file containing an
ext4 filesystem for Firecracker to boot. In order to populate this filesystem
in the file-based block device, Ignite needs to temporarily `mount` the filesystem,
and copy the desired root filesystem source contents in. `mount` requires the UID
to be 0 (root).

We hope to remove this requirement from the Ignite CLI in the future
[#24](https://github.com/weaveworks/ignite/issues/24), [#33](https://github.com/weaveworks/ignite/issues/33).
However, some part of Ignite (although hidden) will always need to execute as root due
to the need to `mount`.

## Q: Can I run Ignite on a Mac?

No. Firecracker requires KVM, as per the above, a feature that is not available on MacOS.
Technically, you could spin up a VM running Linux inside of a Mac, and inside of that Linux
VM, with nested virtualization enabled, run Ignite. However, that might defeat the purpose of
Ignite on Mac in the first place.

## Q: Why is Docker (containers) needed/used?

Docker, currently the only available container runtime usable by Ignite, is used for a couple of reasons:

1. **Running long-lived processes**: At the very early Ignite PoC stage, we tried to run the Firecracker
   process under `systemd`, but this was in many ways suboptimal. Attaching to the serial console, fetching
   logs, and 2. and 3. were very hard to achieve. Also, we'd need to somehow install the Firecracker binary
   on host. Packaging everything in a container, and running the Firecracker process in that container was a
   natural fit.
1. **Sandboxing the Firecracker process**: Firecracker should not be run on host without sandboxing, as per
   their security model.
   Firecracker provides the [jailer](https://github.com/firecracker-microvm/firecracker/blob/master/docs/jailer.md)
   binary to do sandboxing/isolation from the host for the Firecracker process, but a container does this
   equally well, if not better.
1. **Container Networking**: Using containers, we already know what IP to give the VM. We can integrate with
   e.g. the default docker bridge, docker's `libnetwork` in general, or [CNI](https://github.com/containernetworking/cni).
   This reduces the amount of scope and work needed by Ignite, and keeps our implementation lean. It also directly
   makes Ignite usable alongside normal containers, e.g. on a host running Kubernetes Pods.
1. **OCI compliant operations**: Using an existing container runtime, we do not need to implement everything
   from the OCI spec ourselves. Instead, we can re-use functionality from the runtime, e.g. `pull`, `create`,
   and `export`.

All in all, we do not want to reinvent the wheel. We reuse what we can from existing proven container tools.

## Q: How does my filesystem in a Docker image end up in a Firecracker VM?

In short, we `pull` an OCI image using the container runtime (Docker for now), `create` a new container using
this image, and finally `export` the rootfs of that created container to a tar file. This tar file is then
extracted into the mount point of an ext4-formatted block device file of the OCI image's size. The kernel
OCI image is similarly copied into the rootfs of the container. Lastly, Ignite modifies some well-known files
like `/etc/hosts` and `/etc/resolv.conf` for the VM to work as you would expect it to.

## Q: How does networking work as there are both containers and VMs?

First, Ignite spawns a container using the runtime. In this container, one Ignite component, `ignite-spawn`, is running.
`ignite-spawn` loops the network interfaces inside of the container, and looks for a valid one to use for the VM.
It removes the IP address from the container, and remembers it for later.

Next, `ignite-spawn` creates a `tap` device which Firecracker will use, and bridges the `tap` device with the existing
`veth` interface created by the container runtime. With these two interfaces bridged, all information routed to the
container, will end up in the VM's `tap` interface.

Lastly, `ignite-spawn` spawns the Firecracker process, which starts the VM. The VM is started with the `ip=dhcp` kernel
argument, which makes the kernel automatically do a DHCP request for an IP. The kernel asks for an IP to use, and 
`ignite-spawn` responds with the IP the container initially had.

## Q: Where does Ignite originate from?

As per the announcement blog post: https://www.weave.works/blog/fire-up-your-vms-with-weave-ignite

> Ignite is a clean room implementation of a project Lucas prototyped while on army service.

> Lucas Käldström (@luxas) is a Kubernetes SIG Lead and Top CNCF Ambassador 2017, and is a longstanding member of the Weaveworks family since graduating from High School (story here). As a young Finnish citizen, Lucas had to do his mandatory Military Service for around a year.

> Naturally for Lucas, he started evangelising Kubernetes within the military, and got assigned programming tasks. Security and resource consumption are critical army concerns, so Lucas and a colleague, Dennis Marttinen (@twelho), decided to experiment with Firecracker, creating an elementary version of Ignite. On leaving the army they were granted permission to work on an open source rewrite, working with Weaveworks.


## Q: I'm using containerd as my runtime, how do I see the running containers and images?

Containerd has the concept of namespaces which is different than Docker.  In order to view containers and images used by ignite you need to pass the namespace parameter to `ctr` using either `--namespace firecracker` or `-n firecracker`.  