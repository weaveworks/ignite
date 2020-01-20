# Kernel Images

These kernel OCI images contain the kernel binary (at `/boot/vmlinux`) and supporting modules (in `/lib/modules`)
for guest VMs ran by Ignite.

## Building the Kernel Images

```console
$ make
```

## Versions

All LTS versions starting from 4.14 and above are supported by the Ignite team.
This means in practice:

- 4.14.x
- 4.19.x
- 5.4.x

The exact patch versions may be found in the [Makefile](Makefile).
The available versions exist in the [stable kernel git tree](https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/refs/).

## Upgrading to a new kernel version

It should be as easy as:

```console
$ make upgrade
```

after you've upgraded the values in the Makefile.

## Kernel Config Parameters we care about

Some options to the kernel are specifically important for making guest software work.

Please see: [config-patches](config-patches) for what kernel configs we've changed.
The base kernel config is the MicroVM-optimized config file from the Firecracker team.
We're storing it in [upstream/config-amd64](upstream/config-amd64). It's available online
at [firecracker/resources](https://github.com/firecracker-microvm/firecracker/tree/master/resources).
