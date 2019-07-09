## ignite

ignite: easily run Firecracker VMs

### Synopsis


Ignite is a containerized Firecracker microVM administration tool.
It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

Administration is divided into three subcommands:
  image       Manage base images for VMs
  kernel      Manage VM kernels
  vm          Manage VMs

Ignite also supports the same commands as the Docker CLI.
Combining an Image and a Kernel gives you a runnable VM.

Example usage:

	$ ignite run centos:7 \
		--cpus 2 \
		--memory 2GB \
		--ssh \
		--name my-vm
	$ ignite images
	$ ignite kernels
	$ ignite ps
	$ ignite logs my-vm
	$ ignite ssh my-vm


### Options

```
  -h, --help                 help for ignite
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite attach](ignite_attach.md)	 - Attach to a running VM
* [ignite completion](ignite_completion.md)	 - Output bash completion for ignite to stdout
* [ignite create](ignite_create.md)	 - Create a new VM without starting it
* [ignite gitops](ignite_gitops.md)	 - Run the GitOps feature of Ignite
* [ignite image](ignite_image.md)	 - Manage base images for VMs
* [ignite inspect](ignite_inspect.md)	 - Inspect an Ignite Object
* [ignite kernel](ignite_kernel.md)	 - Manage VM kernels
* [ignite kill](ignite_kill.md)	 - Kill running VMs
* [ignite logs](ignite_logs.md)	 - Get the logs for a running VM
* [ignite ps](ignite_ps.md)	 - List running VMs
* [ignite rm](ignite_rm.md)	 - Remove VMs
* [ignite rmi](ignite_rmi.md)	 - Remove VM base images
* [ignite rmk](ignite_rmk.md)	 - Remove kernels
* [ignite run](ignite_run.md)	 - Create a new VM and start it
* [ignite ssh](ignite_ssh.md)	 - SSH into a running vm
* [ignite start](ignite_start.md)	 - Start a VM
* [ignite stop](ignite_stop.md)	 - Stop running VMs
* [ignite version](ignite_version.md)	 - Print the version of ignite
* [ignite vm](ignite_vm.md)	 - Manage VMs

