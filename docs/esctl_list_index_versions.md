## esctl list index versions

show index creation version

```
esctl list index versions [index pattern] [flags]
```

### Aliases

```
version
```

### Examples

```
# list all indexes and their versions
esctl list index versions

# list all indexes and their versions for pattern
esctl list index versions watch*

```

### Options

```
  -h, --help   help for versions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list index](esctl_list_index.md)	 - show information about one or more index

