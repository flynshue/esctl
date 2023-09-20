package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

var getClusterCmd = &cobra.Command{
	Use:   "cluster [command]",
	Short: "show cluster info",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
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

func putClusterSettings(body string) ([]byte, error) {
	b := bytes.NewBufferString(body)
	resp, err := client.Cluster.PutSettings(b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func init() {
	getCmd.AddCommand(getRebalanceCmd)
	setCmd.AddCommand(setRebalanceThrottleCmd)
	resetCmd.AddCommand(resetRebalanceThrottleCmd)
}
