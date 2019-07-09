## ignite vm start

Start a VM

### Synopsis


Start the given VM. The VM is matched by prefix based on its ID and name.
If the interactive flag (-i, --interactive) is specified, attach to the
VM after starting.


```
ignite vm start <vm> [flags]
```

### Options

```
  -d, --debug         Debug mode, keep container after VM shutdown
  -h, --help          help for start
  -i, --interactive   Attach to the VM after starting
      --net string    Networking mode to use. Available options are: [cni docker-bridge] (default "docker-bridge")
```

### Options inherited from parent commands

```
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite vm](ignite_vm.md)	 - Manage VMs

