## ignited version

Print the version of ignite

### Synopsis

Print the version of ignite

```
ignited version [flags]
```

### Options

```
  -h, --help            help for version
  -o, --output string   Output format; available options are 'yaml', 'json' and 'short'
```

### Options inherited from parent commands

```
      --id-prefix string        Prefix string for identifiers and names (default "ignite")
      --ignite-config string    Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignited](ignited.md)	 - ignited: run Firecracker VMs declaratively through a manifest directory or Git

