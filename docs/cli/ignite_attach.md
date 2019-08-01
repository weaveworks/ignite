## `ignite attach` - Attach to a running VM

Attach to a running VM

### Synopsis


Connect the current terminal to the running VM's TTY.
To detach from the VM's TTY, type ^P^Q (Ctrl + P + Q).
The given VM is matched by prefix based on its ID and name.


```
ignite attach <vm> [flags]
```

### Options

```
  -h, --help   help for attach
```

### Options inherited from parent commands

```
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite](index) - ignite: easily run Firecracker VMs
