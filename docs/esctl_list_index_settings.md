## esctl list index settings

list indexes with a summary of settings. Includes replicas, shards, ilm policy, ilm rollover alias, and auto expand replicas

```
esctl list index settings [index pattern] [flags]
```

### Aliases

```
config, cfg
```

### Examples

```
# List all indexes with summary of settings
esctl list index settings

# List indexes matching pattern with summary of settings
esctl list index settings .fleet-*

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

