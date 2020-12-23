## ignite run

Create a new VM and start it

### Synopsis


Create and start a new VM immediately. The image (and kernel) is matched by
prefix based on its ID and name. This command accepts all flags used to
create and start a VM. The interactive flag (-i, --interactive) can be
specified to immediately attach to the started VM after creation.

Example usage:
	$ ignite run weaveworks/ignite-ubuntu \
		--interactive \
		--name my-vm \
		--cpus 2 \
		--ssh \
		--memory 2GB \
		--size 10G


```
ignite run <OCI image> [flags]
```

### Options

```
      --config string                     Specify a path to a file with the API resources you want to pass
  -f, --copy-files strings                Copy files/directories from the host to the created VM
      --cpus uint                         VM vCPU count, 1 or even numbers between 1 and 32 (default 1)
  -d, --debug                             Debug mode, keep container after VM shutdown
  -h, --help                              help for run
      --ignore-preflight-checks strings   A list of checks whose errors will be shown as warnings. Example: 'BinaryInPath,Port,ExistingFile'. Value 'all' ignores errors from all checks.
  -i, --interactive                       Attach to the VM after starting
      --kernel-args string                Set the command line for the kernel (default "console=ttyS0 reboot=k panic=1 pci=off ip=dhcp")
  -k, --kernel-image oci-image            Specify an OCI image containing the kernel at /boot/vmlinux and optionally, modules (default weaveworks/ignite-kernel:4.19.125)
  -l, --label stringArray                 Set a label (foo=bar)
      --memory size                       Amount of RAM to allocate for the VM (default 512.0 MB)
  -n, --name string                       Specify the name
      --network-plugin plugin             Network plugin to use. Available options are: [cni docker-bridge] (default cni)
  -p, --ports strings                     Map host ports to VM ports
      --require-name                      Require VM name to be passed, no name generation
      --runtime runtime                   Container runtime to use. Available options are: [docker containerd] (default containerd)
      --sandbox-image oci-image           Specify an OCI image for the VM sandbox (default weaveworks/ignite:dev)
  -s, --size size                         VM filesystem size, for example 5GB or 2048MB (default 4.0 GB)
      --ssh[=<path>]                      Enable SSH for the VM. If <path> is given, it will be imported as the public key. If just '--ssh' is specified, a new keypair will be generated. (default is unset, which disables SSH access to the VM)
  -v, --volumes volume                    Expose block devices from the host inside the VM
```

### Options inherited from parent commands

```
      --id-prefix string       Prefix string for identifiers and names (default "ignite")
      --ignite-config string   Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel     Specify the loglevel for the program (default info)
  -q, --quiet                  The quiet mode allows for machine-parsable output by printing only IDs
```

### SEE ALSO

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

