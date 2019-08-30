# Networking

Ignite uses network plugins to manage VM networking.

The default plugin is [CNI](https://github.com/containernetworking/cni); and the default CNI network is automatically put in
`/etc/cni/net.d/10-ignite.conflist` if `/etc/cni/net.d` is empty. In order to switch to some other CNI plugin,
remove `/etc/cni/net.d/10-ignite.conflist`, and install e.g. [Weave Net](#multi-node-networking-with-weave-net) like below.

The legacy `docker-bridge` networking provider is also supported, but deprecated.

To select the network plugin, use the `--network-plugin` flag for `ignite` and `ignited`:

```console
ignite --network-plugin cni <command>
ignited --network-plugin docker-bridge <command>
```

## Comparison

### The default CNI network

Automatically installed to `/etc/cni/net.d/10-ignite.conflist` unless you have populated `/etc/cni/net.d` with something else.

**Pros:**

- **Multi-node support**: CNI implementations can often route packets between multiple physical hosts. External computers can access the VM's IP.
- **Kubernetes-compatible**: You can use the same overlay networks as you use with Kubernetes, and hence get your VMs on the same network as your containers.
- **Port mapping support**: This mode supports port mappings from the VM to the host

**Cons:**

- **No multi-node support**: The IP is local (in the `172.18.0.0/16` range), and hence other computers can't connect to your VM's IP address.

### An third-party CNI plugin

For example, [Weave Net](#multi-node-networking-with-weave-net), or an other Kubernetes and/or CNI-implementation.

**Pros:**

- **Multi-node support**: CNI implementations can often route packets between multiple physical hosts. External computers can access the VM's IP.
- **Kubernetes-compatible**: You can use the same overlay networks as you use with Kubernetes, and hence get your VMs on the same network as your containers.
- **Port mapping support**: This mode supports port mappings from the VM to the host

**Cons:**

- **More software needed**: There's now one extra piece of software to configure and manage.

### docker-bridge

**Pros:**

- **Quick start**: If you're running `docker`, you can get up and running without installing extra software
- **Port mapping support**: This mode supports port mappings from the VM to the host

**Cons:**

- **docker-dependent**: By design, this mode is can only be used with docker, and is hence not portable across container runtimes.
- **No multi-node support**: The IP is local (in the `172.17.0.0/16` range), and hence other computers can't connect to your VM's IP address.

## Multi-node networking with Weave Net

To use e.g. Ignite together with [Weave Net](https://github.com/weaveworks/weave), run this on all physical machines that
need to connect to the overlay network:

```shell
# Remove Ignite's default CNI network if it exists
rm -rf /etc/cni/net.d/10-ignite.conflist

# This tries to autodetect the primary IP address of this machine
# Ref: https://stackoverflow.com/questions/13322485/how-to-get-the-primary-ip-address-of-the-local-machine-on-linux-and-macos
export PRIMARY_IP_ADDRESS=$(ip -o route get 1.1.1.1 | cut -d' ' -f7)
# A space-separated list of all the peers in the overlay network
export KUBE_PEERS="${PRIMARY_IP_ADDRESS}"
# Start Weave Net in a container
docker run -d \
  --privileged \
  --net host \
  --pid host \
  --restart always \
  -e HOSTNAME="$(hostname)" \
  -e KUBE_PEERS="${KUBE_PEERS}" \
  -v /var/lib/weave:/weavedb \
  -v /opt:/host/opt \
  -v /home:/host/home \
  -v /etc:/host/etc \
  -v /var/lib/dbus:/host/var/lib/dbus \
  -v /lib/modules:/lib/modules \
  -v /run/xtables.lock:/run/xtables.lock \
  --entrypoint /home/weave/launch.sh \
  weaveworks/weave-kube:2.5.2
```

If you're running Kubernetes on the physical machine you want to use for Ignite VMs, it should work out of the box, as
the CNI implementation is most probably already running in a `DaemonSet` on that machine.
