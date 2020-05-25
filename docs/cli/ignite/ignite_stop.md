## ignite stop

Stop running VMs

### Synopsis


Stop one or multiple VMs. The VMs are matched by prefix based on their
ID and name. To stop multiple VMs, chain the matches separated by spaces.
The force flag (-f, --force) kills VMs instead of trying to stop them
gracefully.

The VMs are given a 20 second grace period to shut down before they
will be forcibly killed.


```
ignite stop <vm>... [flags]
```

### Options

```
  -f, --force-kill   Force kill the VM
  -h, --help         help for stop
```

### Options inherited from parent commands

```
      --ignite-config string    Ignite configuration path
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
  -q, --quiet                   The quiet mode allows for machine-parsable output by printing only IDs
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

