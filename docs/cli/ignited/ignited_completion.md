## ignited completion

Output bash completion for ignited to stdout

### Synopsis


In order to start using the auto-completion, run:

	. <(ignited completion)

To configure your bash shell to load completions for each session, run:

	echo '. <(ignited completion)' >> ~/.bashrc


```
ignited completion [flags]
```

### Options

```
  -h, --help   help for completion
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

