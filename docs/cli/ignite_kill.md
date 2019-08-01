## `ignite kill` - Kill running VMs

Kill running VMs

### Synopsis


Kill (force stop) one or multiple VMs. The VMs are matched by prefix based
on their ID and name. To kill multiple VMs, chain the matches separated
by spaces.


```
ignite kill <vm>... [flags]
```

### Options

```
  -h, --help   help for kill
```

### Options inherited from parent commands

```
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite](index) - ignite: easily run Firecracker VMs
