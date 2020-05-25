## ignite kill

Kill running VMs

### Synopsis


Kill (force stop) one or multiple VMs. The VMs are matched by prefix based
on their ID and name. To kill multiple VMs, chain the matches separated
by spaces.


```
ignite kill <vm>... [flags]
```

### Options

```
  -h, --help   help for kill
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

