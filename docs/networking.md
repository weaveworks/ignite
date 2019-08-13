# Networking

Ignite uses network plugins to manage VM networking. The default plugin is `docker-bridge`, which means that the default docker bridge will be
used for the networking setup. The default docker bridge is a local `docker0` interface, giving out local addresses to containers in the `172.17.0.0/16` range.

Via the `cni` plugin Ignite also supports integration with [CNI](https://github.com/containernetworking/cni), the standard networking
interface for Kubernetes and many other cloud-native projects and container runtimes. Note that CNI itself is only an interface, not
an implementation, so if you use this mode you need to install an implementation of this interface. Any implementation that works
with Kubernetes should technically work with Ignite.

To select the network plugin, use the `--network-plugin` flag for `ignite` and `ignited`:
```console
ignite --network-plugin cni <command>
ignited --network-plugin docker-bridge <command>
```

## Comparison

### docker-bridge

**Pros:**

- **Quick start**: If you're running docker, you can get up and running without installing extra software
- **Port mapping support**: This mode supports port mappings from the VM to the host

**Cons:**

- **docker-dependent**: By design, this mode is can only be used with docker, and is hence not portable across container runtimes.
- **No multi-node support**: The IP is local (in the `172.17.0.0/16` range), and hence other computers can't connect to your VM's IP address.

### CNI

**Pros:**

- **Multi-node support**: CNI implementations can often route packets between multiple physical hosts. External computers can access the VM's IP.
- **Kubernetes-compatible**: You can use the same overlay networks as you use with Kubernetes, and hence get your VMs on the same network as your containers.

**Cons:**

- **More software needed**: There's now one extra piece of software to install and manage.
- **No port-mapping support** (yet): For the moment, we haven't implemented port mapping support for this mode.

## Multi-node networking with Weave Net

To use e.g. Ignite together with [Weave Net](https://github.com/weaveworks/weave), run this on all physical machines that
need to connect to the overlay network:

```shell
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
