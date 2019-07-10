# Ignite installation guide

This guide describes the installation and uninstallation process of Ignite.

## System requirements

Ignite runs on any Intel-based `linux/amd64` system with `KVM` support.
AMD support is in alpha (Firecracker limitation).

**Note**: You do **not** need to install any "traditional" QEMU/KVM packages, as long as
there is virtualization support in the processor and kernel it works. 

See [dependencies.md](dependencies.md) for needed dependencies.

## Downloading the binary
Ignite is a currently a single binary application. To install it,
download the binary from the [GitHub releases page](https://github.com/weaveworks/ignite/releases),
save it as `/usr/local/bin/ignite` and make it executable.

To install Ignite from the command line, follow these steps:

```bash
export VERSION=v0.4.0
curl -Lo ignite https://github.com/weaveworks/ignite/releases/download/${VERSION}/ignite
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

## Removing the installation

To completely remove the Ignite installation, execute the following as root:
```bash
# Force-remove all running VMs
ignite ps -q | xargs ignite rm -f
# Remove the data directory
rm -r /var/lib/firecracker
# Remove the Ignite binary
rm /usr/local/bin/ignite
```