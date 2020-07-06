## ignite vm ssh

SSH into a running vm

### Synopsis


SSH into the running VM using the private key created for it during generation.
If no private key was created or wanting to use a different identity file,
use the identity file flag (-i, --identity) to override the used identity file.
The given VM is matched by prefix based on its ID and name.


```
ignite vm ssh <vm> [flags]
```

### Options

```
  -h, --help              help for ssh
  -i, --identity string   Override the vm's default identity file
      --timeout uint32    Timeout waiting for connection in seconds (default 30)
  -t, --tty               Allocate a pseudo-TTY (default true)
```

### Options inherited from parent commands

```
      --ignite-config string    Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
  -q, --quiet                   The quiet mode allows for machine-parsable output by printing only IDs
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignite vm](ignite_vm.md)	 - Manage VMs

