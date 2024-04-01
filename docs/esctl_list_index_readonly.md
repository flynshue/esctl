## esctl list index readonly

show indexes' read_only setting which are enabled (true)

### Synopsis

The disk-based shard allocator may add and remove the index.blocks.read_only_allow_delete block automatically due to flood stage watermark.
Please see https://www.elastic.co/guide/en/elasticsearch/reference/8.11/index-modules-blocks.html#index-block-settings for more details.

You can use 'esctl disable readonly [index/index pattern] to disable readonly on index/index pattern'

```
esctl list index readonly [flags]
```

### Aliases

```
ro
```

### Examples

```
esctl list index readonly
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

* [esctl list index](esctl_list_index.md)	 - show information about one or more index

