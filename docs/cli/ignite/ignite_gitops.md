## ignite gitops

Run the GitOps feature of Ignite

### Synopsis


Run Ignite in GitOps mode watching the given repository. The repository needs
to be publicly cloneable. Ignite will watch for changes in the master branch
by default, overridable with the branch flag (-b, --branch). If any new/changed
VM specification files are found in the repo (in JSON/YAML format), their
configuration will automatically be declaratively applied.

To quit GitOps mode, use (Ctrl + C).


```
ignite gitops <repo-url> [flags]
```

### Options

```
  -b, --branch string   What branch to sync (default "master")
  -h, --help            help for gitops
  -p, --paths strings   What subdirectories to care about. Default the whole repository
```

### Options inherited from parent commands

```
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

