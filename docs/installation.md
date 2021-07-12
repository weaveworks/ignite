# Ignite installation guide

This guide describes the installation and uninstallation process of Ignite.

## System requirements

Ignite runs on most Intel, AMD or ARM (AArch64) based `linux/amd64` systems with `KVM` support.
See the full CPU support table in [dependencies.md](dependencies.md) for more information.

See [cloudprovider.md](cloudprovider.md) for guidance on running Ignite on various cloud providers and suitable instances that you could use.

**NOTE:** You do **not** need to install any "traditional" QEMU/KVM packages, as long as
there is virtualization support in the CPU and kernel it works.

See [dependencies.md](dependencies.md) for needed dependencies.
Look at [arm.md](arm.md) for use on Raspberry Pi and other ARM machines.

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

Alternatively if you want to run this on AWS EC2, use Amazon Linux v2 (https://aws.amazon.com/amazon-linux-2/) and ensure you use a baremetal instance, like the r series like say r5.metal (https://aws.amazon.com/ec2/instance-types/). 

## Installing dependencies

Ignite has a few dependencies (read more in this [doc](dependencies.md)).
Install them on Ubuntu/CentOS like this:  
(Ignite does not depend on docker package version. If you already installed docker-ce, you don't need to replace it to docker.io.)

Ubuntu:

```bash
apt-get update && apt-get install -y --no-install-recommends dmsetup openssh-client git binutils
which containerd || apt-get install -y --no-install-recommends containerd
    # Install containerd if it's not present -- prevents breaking docker-ce installations
```

CentOS:

```bash
yum install -y e2fsprogs openssh-clients git
which containerd || ( yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo && yum install -y containerd.io )
    # Install containerd if it's not present
```

Amazon Linux 2

```bash
yum install -y e2fsprogs openssh-clients git 
which containerd || amazon-linux-extras enable docker && yum install -y containerd
    # Install containerd if it's not present
```

### CNI Plugins

Install the CNI binaries like this:

```shell
export CNI_VERSION=v0.9.1
export ARCH=$([ $(uname -m) = "x86_64" ] && echo amd64 || echo arm64)
sudo mkdir -p /opt/cni/bin
curl -sSL https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-${ARCH}-${CNI_VERSION}.tgz | sudo tar -xz -C /opt/cni/bin
```

Note that the SSH and Git packages are optional; they are only needed if you use
the `ignite ssh` and/or `ignite gitops` commands.

## Downloading the binary

Ignite is a currently a single binary application. To install it,
download the binary from the [GitHub releases page](https://github.com/weaveworks/ignite/releases),
save it as `/usr/local/bin/ignite` and make it executable.

To install Ignite from the command line, follow these steps:

```bash
export VERSION=v0.9.0
export GOARCH=$(go env GOARCH 2>/dev/null || echo "amd64")

for binary in ignite ignited; do
    echo "Installing ${binary}..."
    curl -sfLo ${binary} https://github.com/weaveworks/ignite/releases/download/${VERSION}/${binary}-${GOARCH}
    chmod +x ${binary}
    sudo mv ${binary} /usr/local/bin
done
```

Ignite uses [semantic versioning](https://semver.org), select the version to be installed
by changing the `VERSION` environment variable.

## Verifying the installation

If the installation was successful, the `ignite` CLI and `ignited` daemon
commands should now be available:

### Ignite CLI

```console
$ ignite version
Ignite version: version.Info{Major:"0", Minor:"8", GitVersion:"v0.9.0", GitCommit:"...", GitTreeState:"clean", BuildDate:"...", GoVersion:"...", Compiler:"gc", Platform:"linux/amd64"}
Firecracker version: v0.22.4
Runtime: containerd
```

### Ignited Daemon

Verify the ignited daemon is running and monitoring files in
`/etc/firecracker/manifests/`.
To do this you will need to install the `uuid-runtime` (Ubuntu) package.

The example here comes from a DigitalOcean droplet Region: NCY1, Size: 2GB image:

```bash
# ignited daemon --log-level debug &
# cd /etc/firecracker/manifests/
# touch smoke-test.yml
DEBU[0830] FileWatcher: Registered inotify events [notify.InCloseWrite: "/etc/firecracker/manifests/smoke-test.yml"] for path "/etc/fire cracker/manifests/smoke-test.yml" 
root@serviceubuntu18-<redacted>-desktop-4ku7nqh:/etc/firecracker/manifests# DEBU[0831] FileWatcher: Sending update: MODIFY -> "/etc/firecracker/manifests/smoke-test.yml" 
DEBU[0831] FileWatcher: Dispatched events batch and reset the events cache 
WARN[0831] Ignoring "/etc/firecracker/manifests/smoke-test.yml": unknown API version "" and/or kind ""
```

Now add some VM details (see [Run Ignite VMs Declaratively](https://ignite.readthedocs.io/en/stable/declarative-config#run-ignite-vms-declaratively)
for additional details):

```bash
VMFILE=/etc/firecracker/manifests/smoke-test.yml
tee "$VMFILE" > /dev/null <<EOF
apiVersion: ignite.weave.works/v1alpha4
kind: VM
metadata:
  name: smoke-test
  uid: $(uuidgen)
spec:
  image:
    oci: weaveworks/ignite-ubuntu
  cpus: 2
  diskSize: 3GB
  memory: 800MB
status:
  running: true
EOF
```

The console output should resemble this:

```bash
DEBU[2551] FileWatcher: Registered inotify events [notify.InCloseWrite: "/etc/firecracker/manifests/smoke-test.yml"] for path "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2552] FileWatcher: Sending update: MODIFY -> "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2552] FileWatcher: Dispatched events batch and reset the events cache
DEBU[2552] GenericMappedRawStorage: AddMapping: "vm/d039cbcd-3606-462d-839e-25ac745cd7c5" -> "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2552] SyncStorage: Received update {{CREATE &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}} 0xc0004c7aa0} true
DEBU[2552] SyncStorage: Sent update: {CREATE &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}}
DEBU[2552] FileWatcher: Skipping suspended event MODIFY for path: "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2552] FileWatcher: Registered inotify events [notify.InCloseWrite: "/etc/firecracker/manifests/smoke-test.yml"] for path "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2553] FileWatcher: Sending update: MODIFY -> "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2553] FileWatcher: Dispatched events batch and reset the events cache
DEBU[2553] SyncStorage: Received update {{MODIFY &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}} 0xc0004c7aa0} true
DEBU[2553] FileWatcher: Skipping suspended event MODIFY for path: "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2553] SyncStorage: Sent update: {MODIFY &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}}
DEBU[2553] FileWatcher: Registered inotify events [notify.InCloseWrite: "/etc/firecracker/manifests/smoke-test.yml"] for path "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2554] FileWatcher: Sending update: MODIFY -> "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2554] FileWatcher: Dispatched events batch and reset the events cache
DEBU[2554] SyncStorage: Received update {{MODIFY &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}} 0xc0004c7aa0} true
DEBU[2554] FileWatcher: Registered inotify events [notify.InCloseWrite: "/etc/firecracker/manifests/smoke-test.yml"] for path "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2554] SyncStorage: Sent update: {MODIFY &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}}
DEBU[2554] FileWatcher: Skipping suspended event MODIFY for path: "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2555] FileWatcher: Sending update: MODIFY -> "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2555] FileWatcher: Dispatched events batch and reset the events cache
DEBU[2555] SyncStorage: Received update {{MODIFY &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}} 0xc0004c7aa0} true
DEBU[2555] SyncStorage: Sent update: {MODIFY &TypeMeta{Kind:VM,APIVersion:ignite.weave.works/v1alpha4,}}
DEBU[2555] FileWatcher: Skipping suspended event MODIFY for path: "/etc/firecracker/manifests/smoke-test.yml"
DEBU[2556] FileWatcher: Registered inotify events [notify.InCloseWrite: "/etc/firecracker/manifests/smoke-test.yml"] for path "/etc/firecracker/manifests/smoke-test.yml"
etc.
etc.
```

Now in a new terminal/console:

1. Verify Ignite is tracking the running micro-VM

```bash
# ignite ps
VM ID                   IMAGE                           KERNEL                                  SIZE    CPU
S       MEMORY          CREATED STATUS  IPS             PORTS   NAME
rpfrdqxmffadvn6t        weaveworks/ignite-ubuntu:latest weaveworks/ignite-kernel:5.4.108        1.2 GB  1 4
56.0 MB 43m ago Up 39m  10.61.0.2               smoke-test
```

1. Verify the micro-VM is accessible (username: root, password: root).

```bash
# ignite attach smoke-test
# cat /proc/cpuinfo
```

To detach from the VM's TTY, type ^P^Q (Ctrl + P + Q).

```bash
# ignite stop smoke-test
# ignite rm smoke-test
```

Now you can continue with the [Getting Started Walkthrough](usage.md).

## Removing the installation

To completely remove the Ignite installation, execute the following as root:

```bash
# Force-remove all running VMs
ignite rm -f $(ignite ps -aq)
# Remove the data directory
rm -r /var/lib/firecracker
# Remove the ignite and ignited binaries
rm /usr/local/bin/ignite{,d}
```
