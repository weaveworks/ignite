## ignite vm attach

Attach to a running VM

### Synopsis


Connect the current terminal to the running VM's TTY.
To detach from the VM's TTY, type ^P^Q (Ctrl + P + Q).
The given VM is matched by prefix based on its ID and name.


```
ignite vm attach <vm> [flags]
```

### Options

```
  -h, --help   help for attach
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

* [ignite vm](ignite_vm.md)	 - Manage VMs

