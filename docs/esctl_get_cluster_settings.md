## esctl get cluster settings

Get cluster-wide settings

### Synopsis

Returns cluster-wide settings in json. You can filter responses using --filter-path
This command is useful for when no command exist to pull back the cluster setting you want to find but also don't feel like writing out the full endpoint url using the esc/console command.

```
esctl get cluster settings [flags]
```

### Examples

```
# Get all cluster settings
esctl get cluster settings

# Get cluster settings using filter
esctl get cluster settings --filter-path "**.routing"

# Get cluster settings using multiple filters
esctl get cluster settings --filter-path "**.routing,**.recovery"
	
```

### Options

```
      --filter-path string   takes a comma separated list of filters expressed with the dot notation.
                             See https://www.elastic.co/guide/en/elasticsearch/reference/8.11/common-options.html#common-options-response-filtering.
                             	
  -h, --help                 help for settings
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl get cluster](esctl_get_cluster.md)	 - show cluster info

