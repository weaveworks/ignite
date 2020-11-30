# Run k3s with Ignite VMs

This short guide shows you how to setup Rancher's Lightweight Kubernetes (k3s) distribution with Ignite VMs.

**NOTE:** At the moment, you need to execute all these commands as `root`.

**NOTE:** k3s requires ~2 GiB of memory to function properly. The default 512 MiB of memory is not enough.

The image contains a minimal Ubuntu 20.04 installation with the k3s server running as a daemon on boot.
Thus, spawning a single k3s Ignite VM is sufficient to operate a single-node k3s cluster. To get started:

```bash
ignite run weaveworks/ignite-k3s:latest \
    --cpus 2 \
    --memory 2GB \
    --ssh \
    --name k3s-master
```

We can test the installation by checking the installed k3s version:

```bash
ignite exec k3s-master -- k3s --version

> k3s version v1.19.4+k3s1 (2532c10f)
```

We can also try to administrate the cluster from the host machine, if `kubectl` is installed:

```bash
VM_IP=$(ignite inspect vm k3s-master | jq -r ".status.network.ipAddresses[0]")
ignite exec k3s-master -- cat /etc/rancher/k3s/k3s.yaml > kubeconfig
sed -i'' "s/127.0.0.1/$VM_IP/" kubeconfig
kubectl --kubeconfig=kubeconfig get nodes

# Should list a single Master node, running the same k3s version displayed earlier
```
