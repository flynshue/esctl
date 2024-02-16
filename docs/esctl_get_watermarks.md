## esctl get watermarks

show watermarks when storage marks readonly

### Synopsis

Disk-based shard allocation settings
------------------------------------
Elasticsearch considers the available disk space on a node before deciding whether to allocate new shards to
that node or to actively relocate shards away from that node.
-------------------------------------------------------------------------------------------------------------
* cluster.routing.allocation.disk.watermark.low
	Controls the low watermark for disk usage. It defaults to 85%, meaning that Elasticsearch
	will not allocate shards to nodes that have more than 85% disk used. It can also be set to
	an absolute byte value (like 500mb) to prevent Elasticsearch from allocating shards if
	less than the specified amount of space is available. This setting has no effect on the
	primary shards of newly-created indices but will prevent their replicas from being allocated.
-------------------------------------------------------------------------------------------------------------
* cluster.routing.allocation.disk.watermark.high
	Controls the high watermark. It defaults to 90%, meaning that Elasticsearch will attempt to
	relocate shards away from a node whose disk usage is above 90%. It can also be set to an
	absolute byte value (similarly to the low watermark) to relocate shards away from a node if
	it has less than the specified amount of free space. This setting affects the allocation of
	all shards, whether previously allocated or not.
-------------------------------------------------------------------------------------------------------------
* cluster.routing.allocation.disk.watermark.flood_stage
	Controls the flood stage watermark, which defaults to 95%. Elasticsearch enforces a read-only
	index block (index.blocks.read_only_allow_delete) on every index that has one or more
	shards allocated on the node, and that has at least one disk exceeding the flood stage.
	This setting is a last resort to prevent nodes from running out of disk space. The index
	block is automatically released when the disk utilization falls below the high watermark.
-------------------------------------------------------------------------------------------------------------
*NOTE*
	You cannot mix the usage of percentage values and byte values within these settings. Either
	all values are set to percentage values, or all are set to byte values. This enforcement is so
	that Elasticsearch can validate that the settings are internally consistent, ensuring that the
	low disk threshold is less than the high disk threshold, and the high disk threshold is less
	than the flood stage threshold.
-------------------------------------------------------------------------------------------------------------

Source: https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-cluster.html#disk-based-shard-allocation

```
esctl get watermarks [flags]
```

### Options

```
  -h, --help   help for watermarks
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl get](esctl_get.md)	 - get details for a resource

