## ignite cp

Copy a file into a running vm

### Synopsis


Copy a file from host into running VM.
Uses SCP to SSH into the running VM using the private key created for it during generation.
If no private key was created or wanting to use a different identity file,
use the identity file flag (-i, --identity) to override the used identity file.
The given VM is matched by prefix based on its ID and name.
Use (-r, --recursive) to recursively copy a directory.


```
ignite cp <vm> <source> <dest> [flags]
```

### Options

```
  -h, --help              help for cp
  -i, --identity string   Override the vm's default identity file
  -r, --recursive         Recursively copy entire directories.
  -t, --timeout uint32    Timeout waiting for connection in seconds (default 10)
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

