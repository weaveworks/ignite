## ignite image build

Build a new base image for VMs

### Synopsis


Build a new base image for VMs. The base image is an ext4
block device file, which contains a root filesystem.

"build" can take in either a tarfile or a Docker image as the source.
The Docker image needs to exist on the host system (pulled locally).

If the import kernel flag (-k, --import-kernel) is specified,
/boot/vmlinux is extracted from the image and added to a new
VM kernel object named after the flag.

Example usage:
	$ ignite build my-image.tar
    $ ignite build luxas/ubuntu-base:18.04 \
		--name my-image \
		--import-kernel my-kernel


```
ignite image build <source> [flags]
```

### Options

```
  -h, --help                   help for build
  -k, --import-kernel string   Import a new kernel from /boot/vmlinux in the image with the specified name
  -n, --name string            Specify the name
```

### Options inherited from parent commands

```
  -q, --quiet   The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite image](ignite_image.md) - Manage VM base images
