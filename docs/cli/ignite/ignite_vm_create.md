## ignite vm create

Create a new VM without starting it

### Synopsis


Create a new VM by combining the given image with a kernel. If no
kernel is given using the kernel flag (-k, --kernel-image), use the
default kernel (weaveworks/ignite-kernel:4.19.125).

Various configuration options can be set during creation by using
the flags for this command.

If the name flag (-n, --name) is not specified,
the VM is given a random name. Using the copy files
flag (-f, --copy-files), additional files/directories
can be added to the VM during creation with the syntax
/host/path:/vm/path.

Example usage:
	$ ignite create weaveworks/ignite-ubuntu \
		--name my-vm \
		--cpus 2 \
		--ssh \
		--memory 2GB \
		--size 6GB


```
ignite vm create <OCI image> [flags]
```

### Options

```
      --config string             Specify a path to a file with the API resources you want to pass
  -f, --copy-files strings        Copy files/directories from the host to the created VM
      --cpus uint                 VM vCPU count, 1 or even numbers between 1 and 32 (default 1)
  -h, --help                      help for create
      --kernel-args string        Set the command line for the kernel (default "console=ttyS0 reboot=k panic=1 pci=off ip=dhcp")
  -k, --kernel-image oci-image    Specify an OCI image containing the kernel at /boot/vmlinux and optionally, modules (default weaveworks/ignite-kernel:4.19.125)
  -l, --label stringArray         Set a label (foo=bar)
      --memory size               Amount of RAM to allocate for the VM (default 512.0 MB)
  -n, --name string               Specify the name
      --network-plugin plugin     Network plugin to use. Available options are: [cni docker-bridge] (default cni)
  -p, --ports strings             Map host ports to VM ports
      --require-name              Require VM name to be passed, no name generation
      --runtime runtime           Container runtime to use. Available options are: [docker containerd] (default containerd)
      --sandbox-image oci-image   Specify an OCI image for the VM sandbox (default weaveworks/ignite:dev)
  -s, --size size                 VM filesystem size, for example 5GB or 2048MB (default 4.0 GB)
      --ssh[=<path>]              Enable SSH for the VM. If <path> is given, it will be imported as the public key. If just '--ssh' is specified, a new keypair will be generated. (default is unset, which disables SSH access to the VM)
  -v, --volumes volume            Expose block devices from the host inside the VM
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

