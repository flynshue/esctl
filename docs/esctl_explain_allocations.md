## esctl explain allocations

Provides an explanation for a shard's current allocation. Typically used to explain unassigned shards.

### Synopsis

Elasticsearch retrieves an allocation explanation for an arbitrary unassigned primary or replica shard.

```
esctl explain allocations [flags]
```

### Aliases

```
alloc
```

### Examples

```
# explain shard allocations
esctl explain allocations

# using cmd alias
esctl explain alloc
	
```

### Options

```
  -h, --help   help for allocations
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl explain](esctl_explain.md)	 - Provides explanation for cluster settings/allocations on resources

