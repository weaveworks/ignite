# Ignite installation guide

This guide describes the installation and uninstallation process of Ignite.

## System requirements

Ignite runs on any Intel-based `linux/amd64` system with `KVM` support.
AMD support is in alpha (Firecracker limitation).

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

#### Cloud Provider KVM supported instances

If you intend to use a cloud provider to test Ignite, you can follow the instructions below.

##### Amazon Web Services

Amazon EC2 [bare metal instances](https://aws.amazon.com/about-aws/whats-new/2018/05/announcing-general-availability-of-amazon-ec2-bare-metal-instances/) provide direct access to the  Intel® Xeon® Scalable processor and memory resources of the underlying server. These instances are ideal for workloads that require access to the hardware feature set (such as Intel® VT-x), for applications that need to run in non-virtualized environments for licensing or support requirements, or for customers who wish to use their own hypervisor.

Here's a list of instances with KVM support, with pricing (as of July 2019), to help you test Ignite. All the instances listed below are EBS-optimized, with 25 Gigabit available network performance and IPv6 support.

| Family | Type | Pricing (US-West-2) per On Demand Linux Instance Hr | vCPUs | Memory (GiB) | Instance Storage (GB) | 
| ---- | ---- | :----: | :----: | :----: | ---- | 
|Compute optimized | c5.metal | $4.08 | 96 |192 |EBS only | 
| General purpose | m5.metal | $4.608 | 96 | 384 | EBS only |
| General purpose |  m5d.metal | $5.424 | 96 | 384  |4 x 900 (SSD) |
|Memory optimized| r5.metal| $6.048 |96 |768| EBS only| 
|Memory optimized| r5d.metal| $6.912 | 96 |768 |4 x 900 (SSD)| 
|Memory optimized| z1d.metal| $4.464 | 48 |384 |2 x 900 (SSD)|
|Storage optimized| i3.metal| $4.992 | 72 | 512 | 8 x 1900 (SSD) |


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
