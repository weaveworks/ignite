## ignite preflight

Checks dependencies are fullfilled

### Synopsis


Run preflight checkers to verify all the required dependencies are present


```
ignite preflight [flags]
```

### Options

```
  -h, --help                       help for preflight
      --ignore-preflight strings   ignore listed preflights
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

