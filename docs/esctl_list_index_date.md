## esctl list index date

list all indexes with their creation date

```
esctl list index date [index Pattern] [flags]
```

### Examples

```
# List indexes and their creation date that match index pattern .fleet*
esctl list index date .fleet*
	
```

### Options

```
  -h, --help    help for date
      --local   display index creation timestamps in local time instead of UTC. Default is false.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list index](esctl_list_index.md)	 - show information about one or more index

