## esctl list snapshots

list summary of snapshots for stored in one or more repositories

```
esctl list snapshots [repo name] [flags]
```

### Aliases

```
snapshot, snap, snaps
```

### Examples

```
# list all snapshots for all repositories
esctl list snapshots

# list all snapshots stored under repository
esctl list snapshots test-elastic-fs
	
```

### Options

```
  -h, --help    help for snapshots
      --local   display snapshot start/end times in local time
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list](esctl_list.md)	 - list information for resource/s
* [esctl list snapshots repository](esctl_list_snapshots_repository.md)	 - list snapshot repositories

