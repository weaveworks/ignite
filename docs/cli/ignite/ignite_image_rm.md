## ignite image rm

Remove VM base images

### Synopsis


Remove one or multiple VM base images. Images are matched by prefix based on
their ID and name. To remove multiple images, chain the matches separated by spaces.
The force flag (-f, --force) kills and removes any running VMs using the image.


```
ignite image rm <image>... [flags]
```

### Options

```
  -f, --force   Force this operation. Warning, use of this mode may have unintended consequences.
  -h, --help    help for rm
```

### Options inherited from parent commands

```
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default docker-bridge)
  -q, --quiet                   The quiet mode allows for machine-parsable output by printing only IDs
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default docker)
```

### SEE ALSO

* [ignite image](ignite_image.md)	 - Manage base images for VMs

