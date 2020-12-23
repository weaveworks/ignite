## ignite completion

Output bash completion for ignite to stdout

### Synopsis


In order to start using the auto-completion, run:

	. <(ignite completion)

To configure your bash shell to load completions for each session, run:

	echo '. <(ignite completion)' >> ~/.bashrc


```
ignite completion [flags]
```

### Options

```
  -h, --help   help for completion
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

