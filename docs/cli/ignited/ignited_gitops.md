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
  -b, --branch string           What branch to sync (default "master")
  -h, --help                    help for gitops
      --hosts-file string       What known_hosts file to use for remote verification (default "~/.ssh/known_hosts")
      --https-password string   What password/access token to use when authenticating with Git over HTTPS
      --https-username string   What username to use when authenticating with Git over HTTPS
      --identity-file string    What SSH identity file to use for pushing
      --interval duration       Sync interval for pushing to and pulling from the remote (default 30s)
      --timeout duration        Git operation (clone, push, pull) timeout (default 1m0s)
```

### Options inherited from parent commands

```
      --id-prefix string        Prefix string for identifiers and names (default "ignite")
      --ignite-config string    Ignite configuration path; refer to the 'Ignite Configuration' docs for more details
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignited](ignited.md)	 - ignited: run Firecracker VMs declaratively through a manifest directory or Git

