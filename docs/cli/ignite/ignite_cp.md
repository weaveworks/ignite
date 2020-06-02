## ignite cp

Copy files/folders between a running vm and the local filesystem

### Synopsis


Copy a file between host and a running VM.
Creates an SFTP connection to the running VM using the private key created for
it during generation, and transfers files between the host and VM. If no
private key was created or wanting to use a different identity file, use the
identity file flag (-i, --identity) to override the used identity file.

Example usage:
	$ ignite cp localfile.txt my-vm:remotefile.txt
	$ ignite cp my-vm:remotefile.txt localfile.txt


```
ignite cp <source> <dest> [flags]
```

### Options

```
  -h, --help              help for cp
  -i, --identity string   Override the vm's default identity file
      --timeout uint32    Timeout waiting for connection in seconds (default 30)
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

