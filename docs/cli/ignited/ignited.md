## ignited

ignited: run Firecracker VMs declaratively through a manifest directory or Git

### Synopsis


Ignite is a containerized Firecracker microVM administration tool.
It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

TODO: ignited documentation


### Options

```
  -h, --help                    help for ignited
      --id-prefix string        Prefix string for identifiers and names (default "ignite")
      --ignite-config string    Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignited completion](ignited_completion.md)	 - Output bash completion for ignited to stdout
* [ignited daemon](ignited_daemon.md)	 - Operates in daemon mode and watches /etc/firecracker/manifests for VM specifications to run.
* [ignited gitops](ignited_gitops.md)	 - Run the GitOps feature of Ignite
* [ignited version](ignited_version.md)	 - Print the version of ignite

