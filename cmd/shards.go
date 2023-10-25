package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	units "github.com/docker/go-units"
	"github.com/spf13/cobra"
)

type catShardsJsonResp struct {
	Index  string `json:"index"`
	Shard  string `json:"shard"`
	State  string `json:"state"`
	PriRep string `json:"prirep"`
	Docs   string `json:"docs"`
	Store  string `json:"store"`
	IP     string `json:"ip"`
	Node   string `json:"node"`
}

var (
	shardSort string
	nodeName  string
	bigger    string
)

var listShardsCmd = &cobra.Command{
	Use:     "shards [index pattern]",
	Aliases: []string{"shard"},
	Short:   "show information about one or more shard",
	Example: `# List all shards for every node
esctl list shards

# List all shards for specific node
esctl list shards --node es-data-03
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idxPattern := []string{"*"}
		if len(args) != 0 {
			idxPattern = args
		}
		if nodeName == "" {
			return listShards(idxPattern)
		}
		// switch {
		// case bigger != "":
		// 	return listShardsNodeBigger(nodeName, bigger)
		// default:
		// 	shards, err := listShardsForNode(nodeName)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	printShards(shards)
		// }
		shards, err := listShardsForNode(nodeName, idxPattern)
		if err != nil {
			return err
		}
		printShards(shards)

		return nil
	},
}

var listShardCountCmd = &cobra.Command{
	Use:   "count",
	Short: "List shard count for each node",
	Long: `List shard count for each node

A good rule-of-thumb is to ensure you keep the number of shards per node below 20 per GB heap it 
has configured. A node with a 30GB heap should therefore have a maximum of 600 shards, but the 
further below this limit you can keep it the better. This will generally help the cluster 
stay in good health.

Source: https://www.elastic.co/blog/how-many-shards-should-i-have-in-my-elasticsearch-cluster
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listShardCount()
	},
}

var getShardsCmd = &cobra.Command{
	Use:     "shards [command]",
	Aliases: []string{"shard"},
	Short:   "show information about one or more shard",
}

var getShardAllocationsCmd = &cobra.Command{
	Use:     "allocations",
	Aliases: []string{"alloc"},
	Short:   "Get shard routing allocation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getShardAllocations()
	},
}

var disableShardCmd = &cobra.Command{
	Use:     "shards [command]",
	Aliases: []string{"shard"},
}

var disableShardAllocationsCmd = &cobra.Command{
	Use:     "allocations",
	Aliases: []string{"alloc"},
	Short:   "Disable shard routing allocations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return disableShardAllocations()
	},
}

var enableShardCmd = &cobra.Command{
	Use:     "shards [command]",
	Aliases: []string{"shard"},
}

var enableShardAllocationsCmd = &cobra.Command{
	Use:     "allocations",
	Aliases: []string{"alloc"},
	Short:   "Enable shard routing allocations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableShardAllocations()
	},
}

var retryShardsCmd = &cobra.Command{
	Use:   "shards",
	Short: "Retry unassigned shards",
	RunE: func(cmd *cobra.Command, args []string) error {
		return retryShards()
	},
}

func getShardAllocations() error {
	b, err := getClusterSettings("**.cluster.routing.allocation.enable")
	if err != nil {
		return err
	}
	settings := map[string]clusterSettings{}
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	for k, v := range settings {
		if v.Cluster.Allocation.Enable == "" {
			continue
		}
		fmt.Printf("%s\ncluster.routing.allocation.enable: %s\n", k, v.Cluster.Allocation.Enable)
	}
	return nil
}

func disableShardAllocations() error {
	b, err := setShardAllocations("none")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func setShardAllocations(enable string) ([]byte, error) {
	body := `{
		"transient": {
		  "cluster.routing.allocation.enable":   "%s"
		}
	   }`
	return putClusterSettings(fmt.Sprintf(body, enable))

}

func enableShardAllocations() error {
	b, err := setShardAllocations("all")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func listShards(idxPattern []string) error {
	b, err := catShards("", idxPattern)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func listShardsNodeBigger(node, size string, idxPattern []string) error {
	resp, err := listShardsForNode(node, idxPattern)
	if err != nil {
		return err
	}
	limit, _ := units.FromHumanSize(size)
	shards := make([]catShardsJsonResp, 0, len(resp))
	for _, r := range resp {
		s, _ := units.FromHumanSize(r.Store)
		if s >= limit {
			shards = append(shards, r)
		}
	}
	printShards(shards)
	return nil
}

func catShards(format string, idxPattern []string) ([]byte, error) {
	resp, err := client.Cat.Shards(client.Cat.Shards.WithHuman(),
		client.Cat.Shards.WithPretty(),
		client.Cat.Shards.WithS(fmt.Sprintf("store:%s,index,shard", shardSort)),
		client.Cat.Shards.WithV(true),
		client.Cat.Shards.WithFormat(format),
		client.Cat.Shards.WithIndex(idxPattern...),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func listShardsForNode(node string, idxPattern []string) ([]catShardsJsonResp, error) {
	b, err := catShards("json", idxPattern)
	if err != nil {
		return nil, err
	}
	resp := []catShardsJsonResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	shards := make([]catShardsJsonResp, 0, len(resp))
	for _, r := range resp {
		if r.Node == node {
			shards = append(shards, r)
		}
	}
	return shards, nil
}

func printShards(shards []catShardsJsonResp) {
	w := newTabWriter()
	fmt.Fprintln(w, "index\t shard\t prirep\t state\t docs\t store\t ip\t node\t")
	for _, shard := range shards {
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t %s\t %s\t\n",
			shard.Index, shard.Shard, shard.PriRep, shard.State, shard.Docs, shard.Store, shard.IP, shard.Node)
	}
	w.Flush()
}

func listShardCount() error {
	b, err := getNodeStats("indices", "nodes.**.name,nodes.**.indices.shard_stats.total_count")
	if err != nil {
		return err
	}
	nodeStats := &nodeStatsResp{}
	if err := json.Unmarshal(b, nodeStats); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "name\t shard_count\t")
	for _, n := range nodeStats.Nodes {
		if strings.Contains(n.Name, "data") {
			fmt.Fprintf(w, "%s\t %d\t\n", n.Name, n.IndexStats.Total)
		}
	}
	w.Flush()
	return nil
}

func retryShards() error {
	resp, err := client.Cluster.Reroute(client.Cluster.Reroute.WithPretty(),
		client.Cluster.Reroute.WithExplain(true),
		client.Cluster.Reroute.WithRetryFailed(true),
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

func init() {
	getCmd.AddCommand(getShardsCmd)
	getShardsCmd.AddCommand(getShardAllocationsCmd)
	listCmd.AddCommand(listShardsCmd)
	listShardsCmd.AddCommand(listShardCountCmd)
	listShardsCmd.Flags().StringVarP(&shardSort, "sort", "s", "desc", "sort shard by size. Valid values are asc or desc. Default is desc.")
	listShardsCmd.Flags().StringVar(&nodeName, "node", "", "filter shards based on node name")
	// listShardsCmd.Flags().StringVar(&bigger, "big", "", "show shards bigger than or equal to size on node, i.e 1gb, 50gb. This only works when supplying --node")
	disableCmd.AddCommand(disableShardCmd)
	disableShardCmd.AddCommand(disableShardAllocationsCmd)
	enableCmd.AddCommand(enableShardCmd)
	enableShardCmd.AddCommand(enableShardAllocationsCmd)
	retryCmd.AddCommand(retryShardsCmd)
}
