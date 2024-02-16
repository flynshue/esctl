package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var disableDestructiveRequiresCmd = &cobra.Command{
	Use:   "destructive-requires",
	Short: "disables destructive_requires_name, wildcards are allowed for deleting indexing",
	RunE: func(cmd *cobra.Command, args []string) error {
		return disableDestructiveRequires()
	},
}

var enableDestructiveRequiresCmd = &cobra.Command{
	Use:   "destructive-requires",
	Short: "enables destructive_requires_name, must specify index name to delete an index. Wildcards are not allowed",
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableDestructiveRequires()
	},
}

var getClusterCmd = &cobra.Command{
	Use:   "cluster [command]",
	Short: "show cluster info",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var getDestructiveRequiresCmd = &cobra.Command{
	Use:   "destructive-requires",
	Short: "destructive_requires_name setting determines if must specify the index name to delete an index or if you can use wildcards.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getDestructiveRequires()
	},
}

var getRebalanceCmd = &cobra.Command{
	Use:     "rebalance-throttle",
	Aliases: []string{"throttle"},
	Short:   "show routing allocations for rebalancing and recoveries",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getClusterRebalance()
	},
}

var getWatermarksCmd = &cobra.Command{
	Use:   "watermarks",
	Short: "show watermarks when storage marks readonly",
	Long: `Disk-based shard allocation settings
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

Source: https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-cluster.html#disk-based-shard-allocation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getClusterWatermarks()
	},
}

var setRebalanceThrottleCmd = &cobra.Command{
	Use:     "rebalance-throttle [size in megabytes]",
	Aliases: []string{"throttle"},
	Short:   "Set bytes per sec routing allocations for rebalancing and recoveries",
	Long: `Set bytes per sec routing allocations for rebalancing and recoveries
size in megabytes: [40|100|250|500|2000|etc.]
NOTE: ...minimum is 40, the max. 2000!...
	`,
	Example: `# Set the rebalance throttle to 250 mb
esctl set rebalance-throttle 250

# same as above, but using the alias cmd
esctl set throttle 250
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply the throttle size")
		}
		size, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if size < 40 || size > 2000 {
			return fmt.Errorf("size in megabytes must be between 40 and 2000")
		}
		return setRebalanceThrottle(size)
	},
}

var resetRebalanceThrottleCmd = &cobra.Command{
	Use:     "rebalance-throttle",
	Aliases: []string{"throttle"},
	Short:   "reset routing allocations for rebalancing, recovery, and throttle to defaults.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return resetRebalanceThrottle()
	},
}

var explainClusterAllocationCmd = &cobra.Command{
	Use:     "allocations",
	Aliases: []string{"alloc"},
	Short:   "Provides an explanation for a shard's current allocation. Typically used to explain unassigned shards.",
	Long:    "Elasticsearch retrieves an allocation explanation for an arbitrary unassigned primary or replica shard.",
	Example: `# explain shard allocations
esctl explain allocations

# using cmd alias
esctl explain alloc
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return explainClusterAllocation()
	},
}

func enableDestructiveRequires() error {
	body := `{
		"persistent": {
			"action.destructive_requires_name" : "true"
		}
	}`
	b, err := putClusterSettings(body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func disableDestructiveRequires() error {
	body := `{
		"persistent": {
			"action.destructive_requires_name" : "false"
		}
	}`
	b, err := putClusterSettings(body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func getClusterRebalance() error {
	b, err := getClusterSettings("**.cluster.routing.allocation,**.indices.recovery.max_bytes_per_sec")
	if err != nil {
		return err
	}
	settings := map[string]clusterSettings{}
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	c := clusterSettings{}
	for _, v := range settings {
		if v.ClusterConcurrentRebalance != "" {
			c.ClusterConcurrentRebalance = v.ClusterConcurrentRebalance
		}
		if v.NodeConcurrentIncomingRecov != "" {
			c.NodeConcurrentIncomingRecov = v.NodeConcurrentIncomingRecov
		}
		if v.NodeConcurrentOutgoingRecov != "" {
			c.NodeConcurrentOutgoingRecov = v.NodeConcurrentOutgoingRecov
		}
		if v.NodeConcurrentRecov != "" {
			c.NodeConcurrentRecov = v.NodeConcurrentRecov
		}
		if v.NodeInitialPriRecov != "" {
			c.NodeInitialPriRecov = v.NodeInitialPriRecov
		}
		if v.Type != "" {
			c.Type = v.Type
		}
		if v.Indices.MaxBps != "" {
			c.Indices.MaxBps = v.Indices.MaxBps
		}
	}
	fmt.Println("cluster.routing.allocation.cluster_concurrent_rebalance: ", c.ClusterConcurrentRebalance)
	fmt.Println("node_concurrent_incoming_recoveries: ", c.NodeConcurrentIncomingRecov)
	fmt.Println("node_concurrent_outgoing_recoveries: ", c.NodeConcurrentOutgoingRecov)
	fmt.Println("node_concurrent_recoveries: ", c.NodeConcurrentRecov)
	fmt.Println("node_initial_primaries_recoveries: ", c.NodeInitialPriRecov)
	fmt.Println("cluster.routing.allocation.type: ", c.Type)
	fmt.Println("indices.recovery.max_bytes_per_sec: ", c.Indices.MaxBps)
	return nil
}

func getClusterSettings(filterPath string) ([]byte, error) {
	resp, err := client.Cluster.GetSettings(client.Cluster.GetSettings.WithIncludeDefaults(true),
		client.Cluster.GetSettings.WithPretty(),
		client.Cluster.GetSettings.WithHuman(),
		client.Cluster.GetSettings.WithFlatSettings(false),
		client.Cluster.GetSettings.WithFilterPath(filterPath),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getDestructiveRequires() error {
	b, err := getClusterSettings("**.action.destructive_requires_name")
	if err != nil {
		return err
	}
	settings := map[string]clusterSettings{}
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	for k, v := range settings {
		if v.Action.Destructive == "" {
			continue
		}
		fmt.Printf("%s\naction.destructive_requires_name: %s\n", k, v.Action.Destructive)
	}
	return nil
}

func setRebalanceThrottle(size int) error {
	body := `{
		"persistent": {
			"indices.recovery.max_bytes_per_sec" : "%dmb"
		}
	}`
	b, err := putClusterSettings(fmt.Sprintf(body, size))
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func resetRebalanceThrottle() error {
	body := `{
		"persistent": {
			"cluster.routing.allocation.cluster_concurrent_rebalance" : null,
			"cluster.routing.allocation.node_concurrent_*" : null,
			"cluster.routing.allocation.node_initial_primaries_recoveries" : null,
			"indices.recovery.max_bytes_per_sec" : null
		}
	}`
	b, err := putClusterSettings(body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func explainClusterAllocation() error {
	resp, err := client.Cluster.AllocationExplain(client.Cluster.AllocationExplain.WithHuman(),
		client.Cluster.AllocationExplain.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func setExcludeNode(name string) error {
	body := fmt.Sprintf(`{
		"persistent" : {
		  "cluster.routing.allocation.exclude._name" : "%s"
		}
	  }`, name)
	b, err := putClusterSettings(body)
	if err != nil {
		log.Println("error putClusterSettings() ", err)
		return err
	}
	fmt.Println(string(b))
	return nil
}

func getExcludedNodes() error {
	b, err := getClusterSettings("**.cluster.routing.allocation.exclude")
	if err != nil {
		return err
	}
	settings := map[string]clusterSettings{}
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	for level, setting := range settings {
		fmt.Printf("%s.cluster.routing.allocation.exclude: %+v\n", level, setting.Exclude)
	}
	return nil
}

func putClusterSettings(body string) ([]byte, error) {
	b := bytes.NewBufferString(body)
	resp, err := client.Cluster.PutSettings(b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getClusterWatermarks() error {
	resp, err := client.Cluster.GetSettings(client.Cluster.GetSettings.WithFlatSettings(true),
		client.Cluster.GetSettings.WithPretty(),
		client.Cluster.GetSettings.WithHuman(),
		client.Cluster.GetSettings.WithIncludeDefaults(true),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	watermarks := map[string]WaterMarksResp{}
	if err := json.Unmarshal(b, &watermarks); err != nil {
		return err
	}
	w := newTabWriter()
	for level, setting := range watermarks {
		empty := WaterMarksResp{}
		if watermarks[level] == empty {
			continue
		}
		fmt.Println(level)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.enable_for_single_data_node:\t %s\t\n", setting.SingleDataNode)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.flood_stage:\t %s\t\n", setting.FloodStage)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.flood_stage.frozen:\t %s\t\n", setting.FloodStageFrozen)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.flood_stage.frozen.max_headroom:\t %s\t\n", setting.FloodStageFrozenMaxHeadRoom)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.flood_stage.max_headroom:\t %s\t\n", setting.FloodStageMaxHeadRoom)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.high:\t %s\t\n", setting.High)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.high.max_headroom:\t %s\t\n", setting.HighMaxHeadRoom)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.low:\t %s\t\n", setting.Low)
		fmt.Fprintf(w, "cluster.routing.allocation.disk.watermark.low.max_headroom:\t %s\t\n\n", setting.LowMaxHeadRoom)
		w.Flush()
	}
	return nil
}

func init() {
	disableCmd.AddCommand(disableDestructiveRequiresCmd)
	enableCmd.AddCommand(enableDestructiveRequiresCmd)
	explainCmd.AddCommand(explainClusterAllocationCmd)
	getCmd.AddCommand(getRebalanceCmd, getDestructiveRequiresCmd, getExcludedNodesCmd, getWatermarksCmd)
	resetCmd.AddCommand(resetRebalanceThrottleCmd)
	setCmd.AddCommand(setRebalanceThrottleCmd, setExcludedNodesCmd)
}
