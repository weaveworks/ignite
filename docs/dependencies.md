# Requirements and dependencies

## Virtualization features

Firecracker by design only supports emulating 4 devices:

- `virtio-block`
- `virtio-net`
- a serial console
- a 1-button keyboard used only to stop the microVM (invoked with `reboot`)

Everything apart from above, is not supported, and out of scope.

## Host Requirements

- A host running Linux 4.14 or newer
- `sysctl net.ipv4.ip_forward=1`
- loaded kernel loop module: `modprobe -v loop`
- Optional: `sysctl net.bridge.bridge-nf-call-iptables=0`, which requires kernel module `br_netfilter`
- One of the following CPUs:

| CPU   | Architecture     | Support level | Notes                                                                                                                                                                         |
|-------|------------------|---------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Intel | x86_64           | Complete      | Requires <a href="https://en.wikipedia.org/wiki/X86_virtualization#Intel_virtualization_(VT-x)">VT-x</a>, most non-Atom 64-bit Intel CPUs since Pentium 4 should be supported |
| AMD   | x86_64           | Alpha         | Requires [AMD-V](https://en.wikipedia.org/wiki/X86_virtualization#AMD_virtualization_.28AMD-V.29), most AMD CPUs since the Athlon 64 "Orleans" should be supported            |
| ARM   | AArch64 (64-bit) | Alpha         | Requires GICv3, see [here](https://github.com/firecracker-microvm/firecracker/issues/1196)                                                                                    |

## Guest Requirements

- A Linux kernel 4.14 or newer
- Kernel config:
  - `CONFIG_VIRTIO_BLK=y` (mandatory)
  - `CONFIG_VIRTIO_NET=y` (mandatory)
  - `CONFIG_KEYBOARD_ATKBD=y` (optional but recommended)
  - `CONFIG_SERIO_I8042=y` (optional but recommended)

## Ignite on-host dependencies

Ignite shells out to a few dependencies on the host.
With time, we aim to eliminate as many of these as possible.

### Container Runtime

- `containerd` for managing the containers Ignite uses (default, preferred)
  - Ubuntu package: `containerd`
  - CentOS package: `containerd.io`
    - From docker's repositories: `yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo`
- `docker` for managing the containers Ignite uses (also installs `containerd` automatically)
  - Ubuntu package: `docker.io`
  - CentOS package: `docker`

### CNI plugins

```shell
export CNI_VERSION=v0.8.2
export ARCH=$([ $(uname -m) = "x86_64" ] && echo amd64 || echo arm64)
curl -sSL https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-${ARCH}-${CNI_VERSION}.tgz | tar -xz -C /opt/cni/bin
```

### Other Binaries

- `mount` & `umount` for mounting and unmounting block devices
  - Ubuntu package: `mount` (installed by default)
  - CentOS package: `util-linux` (installed by default)
- `tar` for extracting files from the docker image onto the filesystem
  - Ubuntu package: `tar` (installed by default)
  - CentOS package: `tar` (installed by default)
- `mkfs.ext4` for formatting a block device with a ext4 filesystem
  - Ubuntu package: `e2fsprogs` (installed by default)
  - CentOS package: `e2fsprogs`
- `e2fsck` & `resize2fs` for cleaning and resizing the ext4 filesystems
  - Ubuntu package: `e2fsprogs` (installed by default)
  - CentOS package: `e2fsprogs`
- `strings` for detecting the kernel version
  - Ubuntu package: `binutils`
  - CentOS package: `binutils` (installed by default)
- `dmsetup` for managing device mapper snapshots and overlays
  - Ubuntu package: `dmsetup`
  - CentOS package: `device-mapper` (installed by default)
- `ssh` for SSH-ing into the VM (optional, for `ignite ssh` only)
  - Ubuntu package: `openssh-client`
  - CentOS package: `openssh-clients`
- `git` for the GitOps mode of Ignite (optional, for `ignite gitops` only)
  - Ubuntu package: `git`
  - CentOS package: `git`
