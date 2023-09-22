## esctl set index auto-expand

Auto-expand the number of replicas based on the number of data nodes in the cluster.

### Synopsis

Auto-expand the number of replicas based on the number of data nodes in the cluster.
Replica range is dash delimited: 0-1, default value is false.
	

```
esctl set index auto-expand [index] [replica range|false] [flags]
```

### Examples

```
# Set auto-expand to 0-1 replicas.
esctl set auto-expand test-logs-0001 0-1

# Disable auto-expand replicas.  Useful if you need to manually set the replicas to 0.
esctl set auto-expand test-logs-0001 false

```

### Options

```
  -h, --help   help for auto-expand
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl set index](esctl_set_index.md)	 - set configuration on index

