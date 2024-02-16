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
	b, err := getClusterSettings("**.cluster.routing.allocation.disk.watermark.low,**.cluster.routing.allocation.disk.watermark.high,**.cluster.routing.allocation.disk.watermark.flood_stage")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
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
