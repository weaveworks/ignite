## ignite kernel import

Import a kernel image from an OCI image

### Synopsis


Import an OCI image as a kernel image for VMs, takes in a Docker image identifier.
This importing is done automatically when the "run" or "create" commands are run.
The import step is essentially a cache for images to be used later when running VMs.


```
ignite kernel import <OCI image> [flags]
```

### Options

```
  -h, --help              help for import
      --runtime runtime   Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### Options inherited from parent commands

```
      --id-prefix string       Prefix string for identifiers and names (default "ignite")
      --ignite-config string   Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel     Specify the loglevel for the program (default info)
  -q, --quiet                  The quiet mode allows for machine-parsable output by printing only IDs
```

### SEE ALSO

* [ignite kernel](ignite_kernel.md)	 - Manage VM kernels

