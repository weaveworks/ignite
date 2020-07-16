# Networking

Ignite uses network plugins to manage VM networking.

The default plugin is [CNI](https://github.com/containernetworking/cni); and the default CNI network is automatically put in
`/etc/cni/net.d/10-ignite.conflist` if `/etc/cni/net.d` is empty. In order to switch to some other CNI plugin,
remove `/etc/cni/net.d/10-ignite.conflist`, and install e.g. [Flannel](#multi-node-networking-with-flannel) like below.

The legacy `docker-bridge` network plugin is also available, but it is deprecated.

To select the network plugin, use the `--network-plugin` flag for `ignite` and `ignited`:

```console
ignite --network-plugin cni <command>
ignited --network-plugin docker-bridge <command>
```

## Comparison

### The default CNI network

Automatically installed to `/etc/cni/net.d/10-ignite.conflist` unless you have populated `/etc/cni/net.d` with something else. Uses the CNI `bridge` plugin.

**Pros:**

- **Kubernetes-compatible**: You can use the same overlay networks as you use with Kubernetes, and hence get your VMs on the same network as your containers.
- **Port mapping support**: This mode supports port mappings from the VM to the host.

**Cons:**

- **No multi-node support**: The default bridge has no logic to communicate with other hosts, local VMs are not discoverable externally. VM IPs are local (in the `10.61.0.0/16` range).

### A third-party CNI plugin

For example [Flannel](#multi-node-networking-with-flannel) or any other Kubernetes and/or CNI-implementation.

**Pros:**

- **Multi-node support**: CNI implementations can often route packets between multiple physical hosts. External computers can access the VM's IP.
- **Kubernetes-compatible**: You can use the same overlay networks as you use with Kubernetes, and hence get your VMs on the same network as your containers.
- **Port mapping support**: This mode supports port mappings from the VM to the host.

**Cons:**

- **More software needed**: There's now one extra piece of software to configure and manage.

**Note:** If you're running Kubernetes on the physical machine you want to use for Ignite VMs, this approach should work
out of the box, as the CNI implementation is most probably already running in a `DaemonSet` on that machine.

### docker-bridge

**Pros:**

- **Quick start**: If you're running `docker`, you can get up and running without installing extra software.
- **Port mapping support**: This mode supports port mappings from the VM to the host.

**Cons:**

- **docker-dependent**: By design, this mode is can only be used with Docker, and is hence not portable across container runtimes.
- **No multi-node support**: The IP is local (in the `172.17.0.0/16` range), and hence other computers can't connect to your VM's IP address.

## Multi-node networking with Flannel

[Flannel](https://github.com/coreos/flannel) is a CNI-compliant layer 3 network fabric. It can be used with Ignite as
a third-party CNI plugin to enable networking across multiple hosts/nodes. To ease the setup process, this repository
provides a helper script at [tools/ignite-flannel.sh](https://github.com/weaveworks/ignite/blob/master/tools/ignite-flannel.sh).

### Configuring the nodes

#### Node 1 (192.168.1.2)

Run `ignite-flannel.sh init` on the first node:

```shell
[node1]$ ./tools/ignite-flannel.sh init
==> Starting ignite-etcd container... 
9a99df0dded30a13a7cd6ec4a04a2038db579ec13c129da53933f3a438474dcd
==> Setting Flannel config:
{
        "Network": "10.50.0.0/16",
        "SubnetLen": 24,
        "SubnetMin": "10.50.10.0",
        "SubnetMax": "10.50.99.0",
        "Backend": {
                "Type": "udp",
                "Port": 8285
        }
}
==> Starting ignite-flannel container... 
25d7f304ade52ad6e5648db8e99cffc78555a1bac01caed5d7401dbf63af2193
==> Setting CNI config...
==> Initialized, now start your Ignite VMs with the CNI network plugin.
```

This will start etcd and Flannel in Docker containers on the first node. Flannel uses etcd to store its configuration.
 
You may now start VMs on this node using `ignite run --network-plugin cni <image>`. To make sure Flannel is active,
verify that the VMs get IP addresses in the `10.50.0.0/16` subnet and that they have internet connectivity.

#### Node 2 (192.168.1.3)

On the second node it's only necessary to run Flannel since the backing etcd is provided by the first node. Check the IP
address or FQDN of the first node and run `ignite-flannel.sh join <first_node_ip_or_fqdn>` on the second node:

```shell
[node2]$ ./tools/ignite-flannel.sh join 192.168.1.2
==> Starting ignite-flannel container... 
01ba5b9a258c5b029ce5412418e998fed1612663d5e8ffe3dcdc33eb5c29dc24
==> Setting CNI config...
==> Complete, now check if joining was successful using 'docker logs ignite-flannel'.
==> If so, go ahead and start your Ignite VMs with the CNI network plugin.
```

Verify that Flannel on the second node has successfully connected to the etcd of the first node using
`docker logs ignite-flannel`:

```shell
[node2]$ docker logs ignite-flannel
I0716 15:18:02.190887       1 main.go:518] Determining IP address of default interface
I0716 15:18:02.192746       1 main.go:531] Using interface with name eth0 and address 192.168.1.3
I0716 15:18:02.192844       1 main.go:548] Defaulting external address to interface address (192.168.1.3)
I0716 15:18:02.193384       1 main.go:246] Created subnet manager: Etcd Local Manager with Previous Subnet: 10.50.31.0/24
I0716 15:18:02.193503       1 main.go:249] Installing signal handlers
I0716 15:18:02.201055       1 main.go:390] Found network config - Backend type: udp
I0716 15:18:02.209864       1 local_manager.go:201] Found previously leased subnet (10.50.31.0/24), reusing
I0716 15:18:02.213133       1 local_manager.go:220] Allocated lease (10.50.31.0/24) to current node (192.168.1.3) 
I0716 15:18:02.214075       1 main.go:305] Setting up masking rules
I0716 15:18:02.239748       1 main.go:313] Changing default FORWARD chain policy to ACCEPT
I0716 15:18:02.239895       1 main.go:321] Wrote subnet file to /run/flannel/subnet.env
I0716 15:18:02.239910       1 main.go:325] Running backend.
I0716 15:18:02.252725       1 main.go:433] Waiting for 22h59m59.937506934s to renew lease
I0716 15:18:02.260713       1 udp_network_amd64.go:100] Watching for new subnet leases
I0716 15:18:02.279792       1 udp_network_amd64.go:195] Subnet added: 10.50.77.0/24
```

If no errors occurred, the overlay network should now be established. Go ahead and start a VM on the second node and
verify Flannel is active by checking the subnet and internet connectivity like for the first node.

If Flannel is throwing errors about the etcd connection:
- Check that you have entered the IP address or FQDN correctly, e.g. verify that you can ping it.
- Make sure that there is no firewall blocking ports `2379/tcp` (etcd) and/or `8285/udp` (Flannel).

At this point you should be able to ping VMs across hosts. Try to ping a VM on the second node from a VM on the first
node and vice versa. If your routes are properly set up (very likely), you can also ping the VMs of a host directly from
another host, so for example the VMs of the second node directly from the first node (outside of a VM).

### Cleanup

To remove all persistent configuration, run `ignite-flannel.sh cleanup` on both hosts.

```shell
[node*]$ ./tools/ignite-flannel.sh cleanup
==> Cleanup complete. To finish removal of non-persistent resources such as generated
==> network interfaces and iptables rules, reboot your system (or remove them by hand).
```

### What about static IPs?

When using CNI, the CNI provider (e.g. Flannel) is responsible for assigning IP addresses to containers (or in this case
the Ignite VMs). Ignite itself only receives an IP from CNI and forwards it to the VM, so it is up to your CNI provider
to persist the IP addresses. See e.g. Flannel's documentation on
[leases and reservations](https://github.com/coreos/flannel/blob/master/Documentation/reservations.md) on how you could
potentially establish this. Right now it is tricky to implement, since Ignite does not support MAC address persistence.

The `ignite-flannel.sh` script is only meant to provide a relatively simple example on how to set up a standalone CNI
network and thus does not have any readily available options to specify static IPs for VMs. That said, it is essentially
just a script to start Flannel and pass it a given configuration, so feel free to take a look at the code in
[tools/ignite-flannel.sh](https://github.com/weaveworks/ignite/blob/master/tools/ignite-flannel.sh) to see how it works
and how you can extend it. Contributions welcome!

#### Static IPs from inside of the VM

If you know the subnet served by your CNI, you could also consider writing the static IP configuration inside the VM
using e.g. [systemd-networkd's configuration](https://www.freedesktop.org/software/systemd/man/systemd.network.html).
This way the VM itself persists the static IP and any routes you may want to add.