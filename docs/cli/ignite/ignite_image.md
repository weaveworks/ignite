## ignite image

Manage base images for VMs

### Synopsis


Groups together functionality for managing VM base images.
Calling this command alone lists all available images.


```
ignite image [flags]
```

### Options

```
  -h, --help   help for image
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
* [ignite image import](ignite_image_import.md)	 - Import a new base image for VMs
* [ignite image ls](ignite_image_ls.md)	 - List available VM base images
* [ignite image rm](ignite_image_rm.md)	 - Remove VM base images

