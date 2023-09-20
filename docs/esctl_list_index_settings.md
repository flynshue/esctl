## esctl list index settings

list indexes with their settings. Includes replicas, shards, ilm policy, ilm rollover alias, and auto expand replicas

```
esctl list index settings [index pattern] [flags]
```

### Aliases

```
config, cfg
```

### Examples

```
# List all indexes with their settings
esctl list index settings

# List indexes matching pattern with their settings
esctl list index .fleet-*

```

### Options

```
  -h, --help   help for settings
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list index](esctl_list_index.md)	 - show information about one or more index

