# Ignite Configuration

Ignite supports global configuration to set the defaults for most of its
configurations. This is defined by the API object `Configuration`. By default,
ignite looks for the configuration file in `/etc/ignite/config.yaml`. It can
also be passed explicitly to the cli using the global flag `--ignite-config`
with path to the configuration file.

Example configuration:

```yaml
apiVersion: ignite.weave.works/v1alpha3
kind: Configuration
metadata:
  name: test-config
spec:
  runtime: containerd
  networkPlugin: cni
  vmDefaults:
    memory: 2GB
    diskSize: 3GB
    cpus: 2
```

This configures ignite to use 2 vCPUs, 2 GB of RAM and 3 GB of disk space with
containerd as the runtime and CNI as the network plugin by default. Any
configuration user provides, in the form of VM config file or flags, overrides
these default configurations.

To check the configuration used by ignite, set the log level to debug.

Ignite using a global configuration file:

```console
$ ignite run weaveworks/ignite-ubuntu --name my-vm --log-level=debug
DEBU[0000] Found default ignite configuration file /etc/ignite/config.yaml
DEBU[0000] Using ignite configuration file /etc/ignite/config.yaml
...
```

Ignite using internal defaults:

```console
$ ignite run weaveworks/ignite-ubuntu --name my-vm --log-level=debug
DEBU[0000] Using ignite default configurations
...
```

Ignite using explicit configuration:

```console
$ ignite run weaveworks/ignite-ubuntu --name my-vm --ignite-config /tmp/ignite-config.yaml --log-level debug 
DEBU[0000] Using ignite configuration file /tmp/ignite-config.yaml
...
```

The full reference format for the `Configuration` kind is as follows:

```yaml
apiVersion: ignite.weave.works/v1alpha3
kind: VM
metadata:
  # Required, the name of the configuration.
  name: [string]
spec:
  # Optional, name of the runtime to use. [containerd or docker].
  runtime: [string]
  # Optional, name of the network plugin to use. [cni or docker-bridge].
  networkPlugin: [string]
  # Optional, default configuration of VM, VM.Spec.
  vmDefaults:
    memory: [size]
    cpus: [uint64]
    ...
```

You can find the full API reference for `Configuration` kind in the
[pkg/apis/](https://github.com/weaveworks/ignite/tree/main/pkg/apis)
subfolder of the project.
