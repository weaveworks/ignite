## Run kubeadm in HA mode with Ignite VMs

This short guide shows you how to setup Kubernetes in HA mode with Ignite VMs.

**Note:** At the moment, you need to execute all these commands as `root`.

**Note:** This guide assumes you have no running containers, in other words, that
the IP of the first docker container that will be run is `172.17.0.2`. You can check
this with `docker run busybox ip addr`.

First set up some files and certificates using `prepare.sh`

```bash
./prepare.sh
```

This will create a kubeadm configuration file, generate the CA cert, give you a kubeconfig file, etc.

### Start the seed master

For the bootstap master, copy over the CA cert and key to use, and the kubeadm config file:

```bash
ignite run weaveworks/ignite-kubeadm:latest \
    --cpus 2 \
    --memory 1GB \
    --ssh \
    --copy-files $(pwd)/run/config.yaml:/kubeadm.yaml \
    --copy-files $(pwd)/run/pki/ca.crt:/etc/kubernetes/pki/ca.crt \
    --copy-files $(pwd)/run/pki/ca.key:/etc/kubernetes/pki/ca.key \
    --name master-0
```

Log into it using `ignite ssh master-0` and an initialize it with `kubeadm`:

```bash
kubeadm init --config /kubeadm.yaml --upload-certs
```

### Join additional masters

Create more master VMs, but copy only the variables we need for joining:

```bash
for i in {1..2}; do
    ignite run weaveworks/ignite-kubeadm:latest \
        --cpus 2 \
        --memory 1GB \
        --ssh \
        --copy-files $(pwd)/run/k8s-vars.sh:/etc/profile.d/02-k8s.sh \
        --name master-${i}
done
```

SSH into each VM with `ignite ssh master-{1,2}`, and join the control plane:

```bash
kubeadm join firekube.luxas.dev:6443 \
    --token ${TOKEN} \
    --discovery-token-ca-cert-hash sha256:${CA_HASH} \
    --certificate-key ${CERT_KEY} \
    --control-plane
```

### Set up a HAProxy loadbalancer locally

```bash
docker run -d -v $(pwd)/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg -p 6443:443 haproxy:alpine
```

### Use kubectl

This will make `kubectl` talk to any of the three masters you've set up, via HAproxy.

```bash
export KUBECONFIG=$(pwd)/run/admin.conf

kubectl get nodes
```

Right now it's expected that the nodes are in state `NotReady`, as CNI networking isn't set up.

#### Install a CNI Network -- Weave Net

We're gonna use [Weave Net](https://github.com/weaveworks/weave).

```bash
kubectl apply -f https://git.io/weave-kube-1.6
```

With this, the nodes should transition into the `Ready` state in a minute or so.

### Watch the cluster heal

Kill the bootstrap master and see the cluster recover:

```bash
ignite rm -f master-0

kubectl get nodes
```

What's happening underneath here is that HAproxy (or any other loadbalancer) notices that
`master-0` is unhealthy, and removes it from the roundrobin list, while etcd also realizes
that one peer is lost, and re-electing a leader amongst the two that are still standing.
When this is done (takes a second or two) the cluster can continue to serve requests as before.
