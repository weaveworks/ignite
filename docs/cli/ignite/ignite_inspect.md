## ignite inspect

Inspect an Ignite Object

### Synopsis


Retrieve information about the given object of the given kind.
The kind can be "image", "kernel" or "vm". The object is matched
by prefix based on its ID and name. Outputs JSON by default, can
be overridden with the output flag (-o, --output).


```
ignite inspect <kind> <object> [flags]
```

### Options

```
  -h, --help              help for inspect
  -o, --output string     Output the object in the specified format (default "json")
  -t, --template string   Format the output using the given Go template
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

