## ignite image import

Import a new base image for VMs

### Synopsis


Import a base image from an OCI image for VMs, takes in a Docker image as the source.
This importing is done automatically when the run or create commands are run. This step
is essentially a cache to be used later when running VMs.


```
ignite image import <OCI image> [flags]
```

### Options

```
  -h, --help   help for import
```

### Options inherited from parent commands

```
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite image](ignite_image.md)	 - Manage base images for VMs

