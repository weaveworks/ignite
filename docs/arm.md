# Using Ignite on ARM64

Ignite has been tested to function on Raspberry Pi 4, as well as arm64 cloud machines.
These instructions should be helpful in using Ignite on other ARM machines such as those in the cloud.

The general takeaway is that Ignite on arm64 Ubuntu should work out of the box, given a kernel that supports KVM.

Ignite depends on [Firecracker's](https://firecracker-microvm.github.io/) ARM 64-bit support.
Ignite and Firecracker are currently built for `arm64`, also known as `aarch64`.

## OS & Kernel

Ignite was tested on a Raspberry Pi 4 with Raspberry Pi OS and Gentoo arm64 OSes.

See the guide at [sakaki-/gentoo-on-rpi-64bit](https://github.com/sakaki-/gentoo-on-rpi-64bit) guide to get started with gentoo, or alternatively use [sakaki-/bcm2711-kernel-bis](https://github.com/sakaki-/bcm2711-kernel-bis) if you want to use Raspberry Pi OS.

## KVM and storage dependencies

For Raspberry Pi OS, you can follow the normal installation procedure as for Ubuntu.

On Gentoo, the packages mentioned via apt in the normal install docs (`dmsetup`,`kvm-ok`, etc.) are available via Gentoo's Portage. In our test setup, we did not need to explicitly install these.

## containerd

For Raspberry Pi OS, you can follow the normal installation procedure as for Ubuntu.

We don't know of a built-in way to install containerd on Gentoo, but you can extract the deb package manually from [Ubuntu's arm64 package](http://ports.ubuntu.com/ubuntu-ports/pool/main/c/containerd/).

## CNI

CNI may be installed with the normal instructions.
Just remember to change the architecture to `arm64`.
The minimum CNI version supported is `v0.8.7`.

## Ignite

Install the arm64 release binaries as documented in the release notes.

## Sandbox / Kernel / OS container images

The Ignite sandbox and kernel images are published using manifest lists, in exactly the same way as for amd64.

For OS images, only Ubuntu ([weaveworks/ignite-ubuntu:latest](https://hub.docker.com/r/weaveworks/ignite-ubuntu/tags)) is tagged to support multiple architectures. For this one however, you will not need to change any of the default image names.

We do publish arch-specific tags if you would like to specify them explicitly or pull them and push them to another registry that doesn't support manifest lists.

If you want to use your own images, make sure that they are built for an armhf or arm64 kernel and userspace. If you'd like to have support for other OSes than Ubuntu on arm64 let us know.
