# Run a set of Ignite VMs with Footloose

If you think of Firecracker as `runc`, you can think of Ignite as `docker`.
Is there then any `docker-compose` thing for Ignite?

Yes, [Footloose](https://github.com/weaveworks/footloose). With Footloose, you can
run containers as Virtual Machines in two modes, either using `docker` or `ignite`.

`docker` isolation with Footloose is good for environments where KVM is not enabled (e.g.
public CI providers, Macs, etc.), and `ignite` isolation is good when you want to run real
VMs.

## Installation

Install [Footloose 0.6.0](https://github.com/weaveworks/footloose/releases/tag/0.6.0) or higher like this:

```shell
export VERSION=0.6.0
curl -sLo footloose https://github.com/weaveworks/footloose/releases/download/${VERSION}/footloose-${VERSION}-linux-x86_64
chmod +x footloose
sudo mv footloose /usr/local/bin/
```

## Get Started

This how you can have Footloose invoke Ignite in a _declaratively_ manner, using a file containing
an API object.

An example file as follows:

```yaml
# footloose.yaml
cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 3
  spec:
    image: weaveworks/ignite-ubuntu:latest
    name: vm%d
    portMappings:
    - containerPort: 22
    # This is by default "docker". However, just set this to "ignite" and it'll work with Ignite :)
    backend: ignite
    # Optional configuration parameters for ignite:
    ignite:
      cpus: 2
      memory: 1GB
      diskSize: 5GB
      kernel: "weaveworks/ignite-kernel:4.19.178"
```

This Footloose API object specifies an Ignite VM with 2 vCPUs, 1GB of RAM, `weaveworks/ignite-kernel:4.19.178` kernel and 5GB of disk.

Given that you have [Footloose](https://github.com/weaveworks/footloose#install) and [Ignite](installation.md) installed, and the above file
created as `footloose.yaml` in the current directory, you can run

```console
$ footloose create
INFO[0000] Docker Image: weaveworks/ignite-ubuntu:latest present locally 
INFO[0000] Creating machine: cluster-vm0 ...
INFO[0002] Creating machine: cluster-vm1 ...
INFO[0004] Creating machine: cluster-vm2 ...
```

SSH into the VM:

```console
$ footloose ssh vm0
Welcome to Ubuntu 18.04.2 LTS (GNU/Linux 4.19.178 x86_64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
 * Support:        https://ubuntu.com/advantage
This system has been minimized by removing packages and content that are
not required on a system that users do not log into.

To restore this content, you can run the 'unminimize' command.

The programs included with the Ubuntu system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Ubuntu comes with ABSOLUTELY NO WARRANTY, to the extent permitted by
applicable law.

root@a07c6f1782c70136:~#
```

Run the following to stop the VMs:

```console
$ footloose stop
INFO[0000] Stopping machine: cluster-vm0 ...
INFO[0002] Stopping machine: cluster-vm1 ...
INFO[0004] Stopping machine: cluster-vm2 ...
```

For more information check out the Footloose [README](https://github.com/weaveworks/footloose#footlooseyaml).
