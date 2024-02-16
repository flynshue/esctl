## esctl set exclude-node

set nodes to be excluded from cluster

```
esctl set exclude-node [node/s] [flags]
```

### Examples

```
# clear excluded nodes
	esctl set exclude-node

	# exclude single node
	esctl set exclude-node es-data-01

	# exclude multiple nodes
	esctl set exclude-node es-data-01 es-data-02
	
```

### Options

```
  -h, --help   help for exclude-node
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl set](esctl_set.md)	 - configure settings on a resource

