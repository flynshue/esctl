## esctl disable readonly

Disable readonly for one or more index

```
esctl disable readonly [index/index pattern] [flags]
```

### Aliases

```
ro
```

### Examples

```
# Disable read only for specific index
esctl disable readonly test-idx-002

# disable read only index pattern
esctl disable readonly test-idx-1*
	
```

### Options

```
  -h, --help   help for readonly
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl disable](esctl_disable.md)	 - disable resource/s

