# Using ignite on ARM

Ignite had been tested to function on Raspberry Pi 3 and 4.
These instructions should be helpful in using ignite on other ARM machines such as those in the cloud.

Ignite depends on [firecracker's](https://firecracker-microvm.github.io/) ARM support.
Ignite and Firecracker are currently built for `arm64`, also known as `aarch64`.

# OS

Ignite was tested on a Raspberry Pi 4 with a Gentoo arm64 base.
See sakaki's [gentoo-on-rpi-64bit](https://github.com/sakaki-/gentoo-on-rpi-64bit) guide to get started.

# KVM and storage dependencies
On Gentoo, the packages mentioned via apt in the normal install docs (`dmsetup`,`kvm-ok`, etc.) are available via Gentoo's Portage. 
In our test setup, we did not need to explicitly install these.

# containerd
You will need to install containerd.
Currently, the containerd project does not publish arm64 binaries on [GitHub Releases](https://github.com/containerd/containerd/releases).
You may build your own using the containerd repo or potentially fetch a build from DockerHub using this [search query](https://hub.docker.com/search?q=containerd&type=image&architecture=arm64); linuxkit has an arm64 Dockerfile and build that is available: ([link](https://hub.docker.com/r/linuxkit/containerd/tags?page=1&name=arm64)).

#### example of extracting containerd from the linuxkit docker image:
We can fetch the binaries from this image. See the linked tags for the most recent SHA.
```shell
docker container create --name ctrd-arm linuxkit/containerd:6ef473a228db6f6ee163f9b9a051102a1552a4ef-arm64 -- echo hi
docker cp ctrd-arm:/usr/bin/ ./linuxkit-containerd-arm
docker rm ctrd-arm

find ./linuxkit-containerd-arm
```

Now you may install these binaries into the appropriate place on your `$PATH` and configure your containerd [`config.toml`](https://github.com/containerd/containerd/blob/master/docs/man/containerd-config.toml.5.md) and [service unit/init](https://github.com/containerd/containerd/blob/master/containerd.service) file.

# Docker
If you choose to use the docker runtime instead of ignite's default support for containerd, you may install Docker's arm64 binaries.
In our test setup, we installed docker with `emerge`.

# CNI
CNI may be installed with the normal instructions.
Just remember to change the architecture to `arm64`.
The minimum CNI version supported is `v0.8.5`.

# Ignite
Install the arm64 release binaries

# Sandbox / Kernel
The ignite sandbox and kernel images are published using manifest lists.
This means that an image like [weaveworks/ignite-ubuntu:latest](https://hub.docker.com/r/weaveworks/ignite-ubuntu/tags) is tagged to support multiple architectures.
You will not need to change any of the default image names.
We do publish arch-specific tags if you would like to specify them explicitly or pull them and push them to another registry that doesn't support manifest lists.

If you want to use your own images, make sure that they are built for an armhf or arm64 kernel and userspace.
