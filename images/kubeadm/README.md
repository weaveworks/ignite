# Run kubeadm in HA mode with Ignite VMs

This short guide shows you how to setup Kubernetes in HA mode with Ignite VMs.

**NOTE:** At the moment, you need to execute all these commands as `root`.

**NOTE:** It is assumed that you start no new VMs between running `prepare.sh` and starting the masters, as IP addresses are computed consecutively

First set up some files and certificates using `prepare.sh` from this directory:

```bash
./prepare.sh
```

This will create a kubeadm configuration file, generate the CA cert, give you a kubeconfig file, etc.

## Start the seed master

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

# Get the IP address of the initial master, for the kubeadm join command below
export MASTER_IP=$(ignite inspect vm master-0 | jq -r ".status.network.ipAddresses[0]")
```

Initialize it with `kubeadm` using `ignite exec`:

```bash
ignite exec master-0 -- kubeadm init --config /kubeadm.yaml --upload-certs
```

## Join additional masters

Create more master VMs, but copy only the variables we need for joining:

```bash
for i in {1..2}; do
    ignite run weaveworks/ignite-kubeadm:latest \
        --cpus 2 \
        --memory 1GB \
        --ssh \
        --name master-${i}
done
```

Use `ignite exec` to join each VM to the control plane:

```bash
source run/k8s-vars.sh
for i in {1..2}; do
    ignite exec master-${i} -- kubeadm join ${MASTER_IP}.xip.io:6443 \
        --token ${TOKEN} \
        --discovery-token-ca-cert-hash sha256:${CA_HASH} \
        --certificate-key ${CERT_KEY} \
        --control-plane
done
```

## Set up a HAProxy loadbalancer locally

```bash
docker run -d -v $(pwd)/run/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg -p 6443:443 haproxy:alpine
```

## Use kubectl

This will make `kubectl` talk to any of the three masters you've set up, via HAproxy.

```bash
export KUBECONFIG=$(pwd)/run/admin.conf

kubectl get nodes
```

Right now it's expected that the nodes are in state `NotReady`, as CNI networking isn't set up.

### Install a CNI Network -- Weave Net

We're going to use [Weave Net](https://github.com/weaveworks/weave).

```bash
kubectl apply -f "https://github.com/weaveworks/weave/releases/download/v2.8.1/weave-daemonset-k8s.yaml"
```

With this, the nodes should transition into the `Ready` state in a minute or so.

## Watch the cluster heal

Kill the bootstrap master and see the cluster recover:

```bash
ignite rm -f master-0

kubectl get nodes
```

What's happening underneath here is that HAproxy (or any other loadbalancer) notices that
`master-0` is unhealthy, and removes it from the roundrobin list. etcd also realizes
that one peer is lost, and re-elects a leader amongst the two that are still standing.
