## esctl delete index

delete index/index pattern

### Synopsis

Starting with Elasticsearch 8.x, by default, the delete index API call does not support wildcards (*) or _all. 
To use wildcards or _all, set the action.destructive_requires_name cluster setting to false.
See https://www.elastic.co/guide/en/elasticsearch/reference/8.10/index-management-settings.html#action-destructive-requires-name

You can use 'esctl disable destructive-requires' to disable this feature and to allow wildcards for deleting index
	

```
esctl delete index [command] [index pattern] [flags]
```

### Aliases

```
idx
```

### Examples

```
# delete specific index
esctl delete index test-logs

# delete multiple index with index pattern
esctl delete index test-logs-*
	
```

### Options

```
      --force   If true, immediately delete without confirmation
  -h, --help    help for index
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl delete](esctl_delete.md)	 - 

