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
 - `mkfs.ext4` for formatting a block device with a ext4 filesystem
 - `docker` for managing the containers ignite uses
 - `e2fsck` & `resize2fs` for cleaning and resizing the ext4 filesystems
 - `dmsetup` for managing devicemapper snapshots and overlays
 - `tar` for extracting files from the docker image onto the filesystem
 - `ssh` for SSH-ing into the VM
