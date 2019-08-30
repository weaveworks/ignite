# How to use Ignite to run VMs

Ignite is a containerized Firecracker microVM administration tool.
It runs and manages virtual machines in separate containers
using [Firecracker](https://firecracker-microvm.github.io/).

This is a quick guide on how to get started with Ignite.
The guide will cover the following topics in order:

- [How to use Ignite to run VMs](#how-to-use-ignite-to-run-vms)
  - [Importing a VM base image](#importing-a-vm-base-image)
  - [Creating a new VM based on the imported image](#creating-a-new-vm-based-on-the-imported-image)
    - [Options for VM generation](#options-for-vm-generation)
  - [Starting a VM](#starting-a-vm)
  - [Inspecting VMs and their resources](#inspecting-vms-and-their-resources)
  - [Accessing a VM](#accessing-a-vm)
    - [Attaching to the TTY](#attaching-to-the-tty)
  - [SSH into the VM](#ssh-into-the-vm)
  - [All in one](#all-in-one)
  - [Stopping a VM](#stopping-a-vm)
  - [Removing a VM](#removing-a-vm)
  - [Removing other resources](#removing-other-resources)

Keep in mind that all Ignite commands require `root` for now.
This will change later.

Here's a [demo][video1] that shows the topics above in action with ignite running on [Amazon EC2 i3.metal instance](https://aws.amazon.com/ec2/instance-types/i3/) which satisfies the /dev/kvm dependency.

[![](http://img.youtube.com/vi/s_O75zt-oBg/0.jpg "ignite running on Amazon EC2 i3.metal instance")][video1]

[video1]: http://www.youtube.com/watch?v=s_O75zt-oBg

Alternatively, you might want to check out the [TGIK](https://github.com/heptio/tgik) deep-dive session from [Joe Beda](https://twitter.com/jbeda) on what Ignite is and how it works:

[![](https://img.youtube.com/vi/aq-wlslJ5MQ/0.jpg "TGIK 082")][video2]

[video2]: https://youtu.be/aq-wlslJ5MQ

## Importing a VM base image

A VM base image (or just `image`) is an OCI container, which contains a filesystem
and has an init system installed. Ignite currently supports importing `images` from
Docker images, for which it has the following command:

```console
# ignite image import <identifier>
```

The identifier can either be the UID of the image in Docker, or its name. If using the
name without specifying a tag, `:latest` is automatically appended. Ignite can also match
a prefix of the given name/UID in any command provided that it's unique, so you can e.g.
enter just the three first letters of a name if they are unique to a single resource.

Go ahead and import the `weaveworks/ignite-ubuntu` Docker image. If it isn't present locally,
Ignite will pull it for you:

```console
# ignite image import weaveworks/ignite-ubuntu
...
INFO[0002] Created image with ID "cae0ac317cca74ba" and name "weaveworks/ignite-ubuntu:latest" 
```

Now the `weaveworks/ignite-ubuntu` image is imported and ready for VM use.

## Creating a new VM based on the imported image

The `images` are read-only references of what every VM based on them should contain.
To create a functional `VM`, Ignite uses `device mapper` to overlay a writable snapshot
on top of the `image`. All changes to the `VM` will be saved in the snapshot.

Let's create a new VM with some options:

```console
# ignite create weaveworks/ignite-ubuntu \
  --name my-vm \
  --cpus 2 \
  --memory 1GB \
  --size 6GB \
  --ssh
...
INFO[0001] Created VM with ID "3c5fa9a18682741f" and name "my-vm" 
```

### Options for VM generation

The previous example tells Ignite to create a `VM` with the name `my-vm` and that it should have
2 CPU cores, 1 GB of RAM, a writable snapshot size of 6 GB and have SSH access enabled.

The snapshot stores a delta compared to the base `image`, so a `--size` of "6GB" enables
storing 6 Gigabytes of data changes (addition or removal).

The `--ssh` flag generates a new private/public key pair
for the `VM` and exports the public key it into the `VM`.
This is used for `ignite ssh <identifier>` later.

All available options can be listed with `ignite create --help`.

## Starting a VM

Starting a created `VM` is very straight forward:

```
# ignite start my-vm
```

The `VM` will be matched by its name or ID (useful if there are similarly named `VMs`).

If no error occured, your `VM` is now running.

## Inspecting VMs and their resources

Ignite currently manages three kinds of resources: `images`, `kernels` and `VMs`.
The `kernels` are quite transparent, and get automatically imported from the docker
image `weaveworks/ignite-kernel:4.19.47` by default (overridable during `create`).

To list the available `kernels`, enter:

```
# ignite kernels
KERNEL ID               NAME                                    CREATED SIZE    VERSION
aefb459546315344        weaveworks/ignite-kernel:4.19.47        61m ago 49.0 MB 4.19.47
```

To list the imported `images`, enter:

```
# ignite images
IMAGE ID                NAME                            CREATED SIZE
cae0ac317cca74ba        weaveworks/ignite-ubuntu:latest 82m ago 268.9 MB
```

And to list the running `VMs`, enter:

```
# ignite ps
VM ID                   IMAGE                           KERNEL                                  CREATED SIZE    CPUS    MEMORY          STATE   IPS             PORTS   NAME
3c5fa9a18682741f        weaveworks/ignite-ubuntu:latest weaveworks/ignite-kernel:4.19.47        63m ago 4.0 GB  2       1.0 GB          Running 172.17.0.3              my-vm
```

To list all `VMs` instead of just running ones, add the `-a` flag to `ps`.

## Accessing a VM

Ignite has two ways to access a CLI in a `VM`, the first option is to attach to the `VM's` TTY
and the other is to SSH into the `VM`.

### Attaching to the TTY

To attach to the running `VM's` TTY, enter:

```console
# ignite attach my-vm
3c5fa9a18682741f
<enter>
Ubuntu 18.04.2 LTS 3c5fa9a18682741f ttyS0

3c5fa9a18682741f login:
```

If nothing is displayed, hit Enter to re-display the login prompt.
Login using the credentials set in the `image` (usually `root` with password `root`).

**To detach** from the TTY, enter the key combination **^P^Q** (Ctrl + P + Q):

```console
root@3c5fa9a18682741f:~# <^P^Q> read escape sequence
$
```

## SSH into the VM

**NOTE:** SSH works only if the `--ssh` flag is specified during `create`. Otherwise there are
no public keys imported into the `VM` and most `images` have password-based root logins
disabled for security reasons.

To SSH into a `VM`, enter:

```console
# ignite ssh my-vm
Welcome to Ubuntu 18.04.2 LTS (GNU/Linux 4.19.47 x86_64)
...
root@3c5fa9a18682741f:~#
```

To exit SSH, just quit the shell process with `exit`.

**NOTE:** Each SSH access spawns its own session, but TTY access
via `attach` is **shared**, every attached user operates the same terminal.

## All in one

Ignite has a shorthand for peforming `image import`, `create`, `start` and possibly also `attach`
all in one command:

```console
# ignite run weaveworks/ignite-ubuntu \
  --name another-vm \
  --cpus 2 \
  --memory 1GB \
  --size 6GB \
  --ssh \
  --interactive
```

This imports the given `image`, creates a new `VM` from it, starts the `VM` and attaches to the `VM's` TTY.

`run` accepts all the flags for `image import`, `create` and `start`. Using the `--interactive`
flag of `start`, an `attach` is performed right after the `VM` has been started.

## Stopping a VM

Ignite `VMs` can be stopped three ways:

1. By running: `# ignite stop my-vm`
2. By running: `# ignite kill my-vm`
3. By issuing the `reboot` command inside the VM

If the `VM's` `kernel` has support for Firecracker's virtual keyboard, `stop` will issue
Ctrl + Alt + Del to gracefully shut down the `VM`. It will wait 20 seconds for Firecracker
to exit, after which the `VM` will be forcibly killed.

`kill` is an alias for `stop -f`, which force-kills the `VM`. **WARNING:** The `VM` is given
no time to close open resources, so this might lead to data loss or filesystem corruption.

Issuing `reboot` inside the `VM` is the recommended way to stop a `VM` that doesn't support
Firecracker's virtual keyboard. By _rebooting_ the `VM` Firecracker _shuts itself down gracefully_.

**NOTE:** Do _not_ enter `shutdown` or `halt` inside the `VM`, this will result in
Firecracker hanging.

## Removing a VM

To remove `VMs` in Ignite, use the following command:

```
# ignite rm my-vm
```

The `VM` needs to not be running for this to succeed. Using the `--force` flag
a running `VM` can also be removed, it will be killed before removal.

## Removing other resources

To remove an `image`, run:

```
# ignite rmi weaveworks/ignite-ubuntu
```

And to remove a `kernel`, run:

```
# ignite rmk weaveworks/ignite-kernel:4.19.47
```

**NOTE:** To fully uninstall all Ignite data, remove the data directory
at `/var/lib/firecracker`. Remember to stop all running `VMs` before doing this.
