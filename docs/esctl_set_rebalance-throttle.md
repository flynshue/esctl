## esctl set rebalance-throttle

Set bytes per sec routing allocations for rebalancing and recoveries

### Synopsis

Set bytes per sec routing allocations for rebalancing and recoveries
size in megabytes: [40|100|250|500|2000|etc.]
NOTE: ...minimum is 40, the max. 2000!...
	

```
esctl set rebalance-throttle [size in megabytes] [flags]
```

### Aliases

```
throttle
```

### Examples

```
# Set the rebalance throttle to 250 mb
esctl set rebalance-throttle 250

# same as above, but using the alias cmd
esctl set throttle 250
	
```

### Options

```
  -h, --help   help for rebalance-throttle
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl set](esctl_set.md)	 - configure settings on a resource

