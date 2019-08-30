## ignite ps

List running VMs

### Synopsis


List all running VMs. By specifying the all flag (-a, --all),
also list VMs that are not currently running.


```
ignite ps [flags]
```

### Options

```
  -a, --all    Show all VMs, not just running ones
  -h, --help   help for ps
```

### Options inherited from parent commands

```
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
  -q, --quiet                   The quiet mode allows for machine-parsable output by printing only IDs
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

