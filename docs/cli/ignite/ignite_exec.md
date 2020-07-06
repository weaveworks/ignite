## ignite exec

execute a command in a running VM

### Synopsis


Execute a command in a running VM using SSH and the private key created for it during generation.
If no private key was created or wanting to use a different identity file,
use the identity file flag (-i, --identity) to override the used identity file.
The given VM is matched by prefix based on its ID and name.


```
ignite exec <vm> <command...> [flags]
```

### Options

```
  -h, --help              help for exec
  -i, --identity string   Override the vm's default identity file
      --timeout uint32    Timeout waiting for connection in seconds (default 30)
  -t, --tty               Allocate a pseudo-TTY
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

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

