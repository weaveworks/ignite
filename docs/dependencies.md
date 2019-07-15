## Requirements and dependencies

### Virtualization features

Firecracker by design only supports emulating 4 devices:
 - `virtio-block`
 - `virtio-net`
 - a serial console
 - a 1-button keyboard used only to stop the microVM (invoked with `reboot`)

Everything apart from above, is not supported, and out of scope.

### Host Requirements

 - A host running Linux 4.14 or newer
 - An Intel or AMD (alpha) CPU
 - `sysctl net.ipv4.ip_forward=1`
 - loaded kernel loop module: `modprobe -v loop`
 - Optional: `sysctl net.bridge.bridge-nf-call-iptables=0`

### Guest Requirements

 - A Linux kernel 4.14 or newer
 - Kernel config:
   - `CONFIG_VIRTIO_BLK=y` (mandatory)
   - `CONFIG_VIRTIO_NET=y` (mandatory)
   - `CONFIG_KEYBOARD_ATKBD=y` (optional but recommended)
   - `CONFIG_SERIO_I8042=y` (optional but recommended)

### Ignite on-host dependencies

Ignite shells out to a few dependencies on the host.
With time, we aim to eliminate as many of these as possible.

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
 - `docker` for managing the containers ignite uses
   - Ubuntu package: `docker.io`
   - CentOS package: `docker`
 - `dmsetup` for managing devicemapper snapshots and overlays
   - Ubuntu package: `dmsetup`
   - CentOS package: `device-mapper` (installed by default)
 - `ssh` for SSH-ing into the VM (optional, for `ignite ssh` only)
   - Ubuntu package: `openssh-client`
   - CentOS package: `openssh-clients`
 - `git` for the GitOps mode of Ignite (optional, for `ignite gitops` only)
   - Ubuntu package: `git`
   - CentOS package: `git`
