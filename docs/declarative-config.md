## Run Ignite VMs declaratively

Flags can be convenient for simple use cases, but have many limitations.
In more advanced use-cases, and to eventually allow GitOps flows, there is
an other way: telling Ignite what to do _declaratively_, using a file containing
an API object.

The first commands to support this feature is `ignite run` and `ignite create`.
An example file as follows: 

```yaml
apiVersion: ignite.weave.works/v1alpha1
kind: VM
metadata:
  name: test-vm
spec:
  image:
    ref: weaveworks/ignite-ubuntu
    type: Docker
  cpus: 2
  diskSize: 3GB
  memory: 800MB
```

This API object specifies a need for 2 vCPUs, 800MB of RAM and 3GB of disk.

We can tell Ignite to make this happen using simply:

```console
$ ignite run --config test-vm.yaml
```
