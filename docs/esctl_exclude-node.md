## esctl exclude-node

Exclude node/s by name (comma-separated). Move shards off of a node prior to shutting it down

```
esctl exclude-node [nodes] [flags]
```

### Examples

```

# Exclude single node
esctl exclude-node es-data-01

# Exclude multiple nodes
esctl exclude-node es-data-01,es-data-02
	
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

* [esctl](esctl.md)	 - CLI tool for interacting with elasticsearch API

