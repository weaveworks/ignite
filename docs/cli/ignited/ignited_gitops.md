## ignited gitops

Run the GitOps feature of Ignite

### Synopsis


Run Ignite in GitOps mode watching the given repository. The repository needs
to be publicly cloneable. Ignite will watch for changes in the master branch
by default, overridable with the branch flag (-b, --branch). If any new/changed
VM specification files are found in the repo (in JSON/YAML format), their
configuration will automatically be declaratively applied.

To quit GitOps mode, use (Ctrl + C).


```
ignited gitops <repo-url> [flags]
```

### Options

```
  -b, --branch string   What branch to sync (default "master")
  -h, --help            help for gitops
  -p, --paths strings   What subdirectories to care about. Default the whole repository
```

### Options inherited from parent commands

```
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default docker-bridge)
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default docker)
```

### SEE ALSO

* [ignited](ignited.md)	 - ignited: run Firecracker VMs declaratively through a manifest directory or Git

