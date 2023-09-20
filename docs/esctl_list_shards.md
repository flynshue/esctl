## esctl list shards

show information about one or more shard

```
esctl list shards [command] [flags]
```

### Aliases

```
shard
```

### Examples

```
# List all shards for every node
esctl list shards

# List all shards for specific node
esctl list shards --node es-data-03
	
```

### Options

```
  -h, --help          help for shards
      --node string   filter shards based on node name
  -s, --sort string   sort shard by size. Valid values are asc or desc. Default is desc. (default "desc")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list](esctl_list.md)	 - list information for resource/s
* [esctl list shards count](esctl_list_shards_count.md)	 - List shard count for each node

