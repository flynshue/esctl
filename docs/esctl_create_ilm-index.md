## esctl create ilm-index

bootstrap initial ilm index and designate it as the write index for the rollover alias specified

### Synopsis

bootstrap initial ilm index and designate it as the write index for the rollover alias specified.
By default, the initial ilm index will be created as <index-prefix-{now/d}-index-suffix> and will use the index prefix as the rollover alias.
	

```
esctl create ilm-index [index prefix] [index suffix] [flags]
```

### Aliases

```
ilm-idx
```

### Examples

```
# bootstrap initial ilm index with name test-filebeat-7d-7.11.2-2023.03.27-000001
esctl create ilm-index test-filebeat-7d-7.11.2 000001	

```

### Options

```
  -h, --help   help for ilm-index
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl create](esctl_create.md)	 - Create resources

