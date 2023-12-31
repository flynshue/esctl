## esctl list shards count

List shard count for each node

### Synopsis

List shard count for each node

A good rule-of-thumb is to ensure you keep the number of shards per node below 20 per GB heap it 
has configured. A node with a 30GB heap should therefore have a maximum of 600 shards, but the 
further below this limit you can keep it the better. This will generally help the cluster 
stay in good health.

Source: https://www.elastic.co/blog/how-many-shards-should-i-have-in-my-elasticsearch-cluster
	

```
esctl list shards count [flags]
```

### Options

```
  -h, --help   help for count
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list shards](esctl_list_shards.md)	 - show information about one or more shard

