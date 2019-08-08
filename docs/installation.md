# Ignite installation guide

This guide describes the installation and uninstallation process of Ignite.

## System requirements

Ignite runs on any Intel-based `linux/amd64` system with `KVM` support.
AMD support is in alpha (Firecracker limitation).

See [cloudprovider.md](cloudprovider.md) for guidance on running Ignite on various cloud providers and suitable instances that you could use.

**Note**: You do **not** need to install any "traditional" QEMU/KVM packages, as long as
there is virtualization support in the CPU and kernel it works. 

See [dependencies.md](dependencies.md) for needed dependencies.

### Checking for KVM support

Please read [dependencies.md](dependencies.md) for the full reference, but if you quickly want
to check if your CPU and kernel supports virtualization, run these commands:

```console
$ lscpu | grep Virtualization
Virtualization:      VT-x

$ lsmod | grep kvm
kvm_intel             200704  0
kvm                   593920  1 kvm_intel
```

Alternatively, on Ubuntu-like systems there's a tool called `kvm-ok` in the `cpu-checker` package.
Check for KVM support using `kvm-ok`:

```console
$ sudo apt-get update && sudo apt-get install -y cpu-checker
...
$ kvm-ok
INFO: /dev/kvm exists
KVM acceleration can be used
```

With this kind of output, you're ready to go!

## Installing dependencies

Ignite has a few dependencies (read more in this [doc](dependencies.md)).
Install them on Ubuntu/CentOS like this:

Ubuntu:

```bash
apt-get update && apt-get install -y --no-install-recommends docker.io dmsetup openssh-client git binutils
```

CentOS:

```bash
yum install -y docker e2fsprogs openssh-clients git
```

Note that the SSH and Git packages are optional; they are only needed if you use
the `ignite ssh` and/or `ignite gitops` commands.

## Downloading the binary

Ignite is a currently a single binary application. To install it,
download the binary from the [GitHub releases page](https://github.com/weaveworks/ignite/releases),
save it as `/usr/local/bin/ignite` and make it executable.

To install Ignite from the command line, follow these steps:

```bash
export VERSION=v0.4.2
curl -fLo ignite https://github.com/weaveworks/ignite/releases/download/${VERSION}/ignite
chmod +x ignite
sudo mv ignite /usr/local/bin
```

Ignite uses [semantic versioning](https://semver.org), select the version to be installed
by changing the `VERSION` environment variable.

## Verifying the installation

If the installation was successful, the `ignite` command should now be available:

```
# ignite version
Ignite version: version.Info{Major:"0", Minor:"4+", GitVersion:"v0.4.0-rc.1", GitCommit:"7e03dc80be894250f9f97ec4d80261fd2fdcd8f4", GitTreeState:"clean", BuildDate:"2019-07-09T19:03:30Z", GoVersion:"go1.12.1", Compiler:"gc", Platform:"linux/amd64"}
Firecracker version: v0.17.0
```

Now you can continue with the [Getting Started Walkthrough](usage.md).

## Removing the installation

To completely remove the Ignite installation, execute the following as root:

```bash
# Force-remove all running VMs
ignite rm -f $(ignite ps -aq)
# Remove the data directory
rm -r /var/lib/firecracker
# Remove the Ignite binary
rm /usr/local/bin/ignite
```
