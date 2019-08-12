# Run Ignite VMs declaratively using Footloose

This how you can have Footloose invoke Ignite in a _declaratively_ manner, using a file containing
an API object.

An example file as follows: 

```yaml
cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 1
  spec:
    image: weaveworks/ignite-kubeadm:latest
    name: master%d
    portMappings:
    - containerPort: 22
    backend: ignite
    ignite:
      cpus: 2
      memory: 4GB
      diskSize: 30GB
      kernel: "weaveworks/ignite-kernel:4.19.47"
      copyFiles:
        "<ABSOLUTE_PATH>/run/pki/ca.crt": "/etc/kubernetes/pki/ca.crt"
        "<ABSOLUTE_PATH>/run/pki/ca.key": "/etc/kubernetes/pki/ca.key"
```

This Footloose API object specifies an Ignite VM with 2 vCPUs, 4GB of RAM, `weaveworks/ignite-kernel:4.19.47` kernel and 30GB of disk.
You can specify the files that need to be copied to the VM using the `copyFiles` property.

We can tell Footloose to fire up Ignite using simply:

```console
$ footloose --config footloose.k8s.yaml create
```

Run the following to stop the vm:
```console
$ footloose --config footloose.k8s.yaml stop
```

Run the following to delete the vm:
```console
$ footloose --config footloose.k8s.yaml delete
```

Here is an [Kubernetes VMs Example](https://github.com/ignite/docs/footloose.k8s.yaml) that configures a 2 VMs cluster.

For more information on [Footloose](https://github.com/weaveworks/footloose/).

