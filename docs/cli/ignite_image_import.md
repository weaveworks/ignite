## ignite image import

Import a new base image for VMs

### Synopsis


Import a new base image for VMs, takes in a Docker image as the source.
The base image is an ext4 block device file, which contains a root filesystem.

If a kernel is found in the image, /boot/vmlinux is extracted from it
and imported to a kernel with the same name.

Example usage:
    $ ignite image import luxas/ubuntu-base:18.04


```
ignite image import <source> [flags]
```

### Options

```
  -h, --help   help for import
```

### Options inherited from parent commands

```
  -q, --quiet   The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite image](ignite_image.md)	 - Manage VM base images

