## ignite inspect

Inspect an Ignite Object

### Synopsis


Retrieve information about the given object of the given kind.
The kind can be "image", "kernel" or "vm". The object is matched
by prefix based on its ID and name. Outputs JSON by default, can
be overridden to YAML with the yaml flag (-y, --yaml).


```
ignite inspect <kind> <object> [flags]
```

### Options

```
  -h, --help            help for inspect
  -o, --output string   Output the object in the specified format (default "json")
```

### Options inherited from parent commands

```
  -q, --quiet   The quiet mode allows for machine-parsable output, by printing only IDs
```

### SEE ALSO

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

