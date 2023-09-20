## esctl get index settings

get full details of settings for index/index pattern

```
esctl get index settings [index pattern] [flags]
```

### Aliases

```
config, cfg
```

### Examples

```
# Get index settings details for specific index
esctl get index settings .fleet-file-data-agent-000001

# Get index settings details for index pattern
esctl get index settings .fleet-*

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

* [esctl get index](esctl_get_index.md)	 - get detailed information about one or more index

