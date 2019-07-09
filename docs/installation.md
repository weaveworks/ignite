# Ignite installation guide

This guide describes the installation and uninstallation process of Ignite.

## System requirements

Ignite runs on any Intel-based `linux/amd64` system with `KVM` enabled.
AMD support is in alpha (Firecracker limitation).

See the [requirements](REQUIREMENTS.md) for needed dependencies.

## Downloading the binary
Ignite is a currently a single binary application. To install it,
download the binary from the [GitHub releases page](https://github.com/weaveworks/ignite/releases),
save it as `/usr/local/bin/ignite` and make it executable.

To install Ignite from the command line, execute the following in a `root` shell
(or use `sudo` for `curl` and `chmod`):
```bash
export VERSION=0.3.0
curl -Lo /usr/local/bin/ignite https://github.com/weaveworks/ignite/releases/download/v${VERSION}/ignite
chmod +x /usr/local/bin/ignite
```

Ignite uses [semantic versioning](https://semver.org), select the version to be installed
by changing the `VERSION` environment variable.

## Verifying the installation

If the installation was successful, the `ignite` command should now be available:
```
# ignite version
Ignite version: version.Info{Major:"0", Minor:"3", GitVersion:"v0.3.0", GitCommit:"9db63f66c8a38c83212d618f6d0d6995b79e07bf", GitTreeState:"clean", BuildDate:"2019-06-18T13:30:59Z", GoVersion:"go1.12.1", Compiler:"gc", Platform:"linux/amd64"}
Firecracker version: v0.16.0
```

## Removing the installation

**NOTE:** Make sure no virtual machines are running before executing this step.

To completely remove the Ignite installation, execute the following as root:
```bash
rm -r /var/lib/firecracker
rm /usr/local/bin/ignite
```