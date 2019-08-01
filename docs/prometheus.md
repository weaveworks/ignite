# Monitor Ignite with Prometheus

Every Ignite VM container exposes [Prometheus](https://prometheus.io) metrics about itself.

These metrics are available with the following command:

```bash
VM_NAME="my-vm"
VM_ID=$(ignite inspect vm ${VM_NAME} | jq -r .metadata.uid)
curl --unix-socket /var/lib/firecracker/vm/${VM_ID}/prometheus.sock http:/metrics
```

This will report metrics for the `ignite-spawn` component, managing the Firecracker daemon inside of the container.
If you want to see how much overhead `ignite-spawn` and `firecracker` have combined for running a VM, you can 
check it with `docker stats`:

```console
$ VM_NAME="my-vm"
$ VM_ID=$(ignite inspect vm ${VM_NAME} | jq -r .metadata.uid)
$ docker stats ignite-${VM_ID}
CONTAINER ID        NAME                      CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
29259f616d9c        ignite-cc82b4424244b3e4   3.12%               9.145MiB / 15.52GiB   0.06%               5.35kB / 2.5kB      324MB / 4.1kB       15
$ docker top ignite-${VM_ID}
UID                 PID                 PPID                C                   STIME               TTY                 TIME                CMD
root                28693               28666               0                   14:11               pts/0               00:00:00            /usr/local/bin/ignite-spawn cc82b4424244b3e4
root                28785               28693               1                   14:11               pts/0               00:00:01            firecracker --api-sock /tmp/firecracker.sock
```
