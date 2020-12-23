## ignite rmk

Remove kernels

### Synopsis


Remove one or multiple VM kernels. Kernels are matched by prefix based on their
ID and name. To remove multiple kernels, chain the matches separated by spaces.
The force flag (-f, --force) kills and removes any running VMs using the kernel.


```
ignite rmk <kernel> [flags]
```

### Options

```
  -f, --force   Force this operation. Warning, use of this mode may have unintended consequences.
  -h, --help    help for rmk
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

