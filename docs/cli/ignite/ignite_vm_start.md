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
  -d, --debug                             Debug mode, keep container after VM shutdown
  -h, --help                              help for start
      --ignore-preflight-checks strings   A list of checks whose errors will be shown as warnings. Example: 'BinaryInPath,Port,ExistingFile'. Value 'all' ignores errors from all checks.
  -i, --interactive                       Attach to the VM after starting
      --network-plugin plugin             Network plugin to use. Available options are: [cni docker-bridge] (default cni)
      --runtime runtime                   Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### Options inherited from parent commands

```
      --id-prefix string       Prefix string for identifiers and names (default "ignite")
      --ignite-config string   Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel     Specify the loglevel for the program (default info)
  -q, --quiet                  The quiet mode allows for machine-parsable output by printing only IDs
```

### SEE ALSO

* [ignite vm](ignite_vm.md)	 - Manage VMs

