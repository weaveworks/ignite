## ignite vm stop

Stop running VMs

### Synopsis


Stop one or multiple VMs. The VMs are matched by prefix based on their
ID and name. To stop multiple VMs, chain the matches separated by spaces.
The force flag (-f, --force) kills VMs instead of trying to stop them
gracefully.

The VMs are given a 20 second grace period to shut down before they
will be forcibly killed.


```
ignite vm stop <vm>... [flags]
```

### Options

```
  -f, --force-kill   Force kill the VM
  -h, --help         help for stop
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

