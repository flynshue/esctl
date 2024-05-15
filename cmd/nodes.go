/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

var clearExcludedNodesCmd = &cobra.Command{
	Use:   "exclude-node",
	Short: "Add nodes that have been excluded from cluster back",
	RunE: func(cmd *cobra.Command, args []string) error {
		return setExcludeNode("")
	},
}

var getExcludedNodesCmd = &cobra.Command{
	Use:   "exclude-node",
	Short: "list nodes that have been excluded from cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getExcludedNodes()
	},
}

var setExcludedNodesCmd = &cobra.Command{
	Use:   "exclude-node [node/s]",
	Short: "set nodes to be excluded from cluster",
	Example: `
	# exclude single node
	esctl set exclude-node es-data-01

	# exclude multiple nodes
	esctl set exclude-node es-data-01 es-data-02
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply node name to exclude")
		}
		return setExcludeNode(strings.Join(args, ","))
	},
}

// nodesCmd represents the nodes command
var listNodesCmd = &cobra.Command{
	Use:     "nodes [command]",
	Aliases: []string{"node"},
	Short:   "show information about one or more node",
}

var nodeStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "list ES nodes with usage statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodes()
	},
}

var nodeSuffixesCmd = &cobra.Command{
	Use:   "suffix",
	Short: "list ES nodes name suffixes",
	RunE: func(cmd *cobra.Command, args []string) error {
		suffixes, err := listNodeNameSuffix()
		if err != nil {
			return err
		}
		fmt.Printf("valid data node suffixes: %s\n", strings.Join(suffixes, ", "))
		return nil
	},
}

var nodeStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "list ES nodes HDD usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodeStorage()
	},
}

var nodeFSDetailsCmd = &cobra.Command{
	Use:     "filesystem",
	Aliases: []string{"fs"},
	Short:   "list ES nodes filesystem details",
	Example: `# basic usage
esctl list nodes filesystem

# using alias
esctl list nodes fs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodeFSDetails()
	},
}

var nodeVersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"versions"},
	Short:   "list ES nodes version",
	Example: `# basic usage
esctl list nodes version

# using alias
esctl list nodes versions
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodesVersion()
	},
}

func listNodes() error {
	b, err := catNodes("ip,name,heap.percent,ram.percent,cpu,load_1m,load_5m,load_15m,node.role,master,name,disk.total,disk.used,disk.avail,disk.used_percent",
		"name:asc",
	)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	suffixes, err := listNodeNameSuffix()
	if err != nil {
		return err
	}
	fmt.Printf("valid data node suffixes: %s\n", strings.Join(suffixes, ", "))
	return nil
}

func listNodeNameSuffix() ([]string, error) {
	nodes, err := listNodesInfo()
	if err != nil {
		return nil, err
	}
	nodeSuffixes := make([]string, 0, len(nodes.Nodes))
	for _, node := range nodes.Nodes {
		if strings.Contains(node.Name, "data") {
			s := strings.SplitAfter(node.Name, "data-")
			nodeSuffixes = append(nodeSuffixes, s[1])
		}
	}
	return nodeSuffixes, nil
}

func listNodeStorage() error {
	b, err := catNodes("ip,node.role,master,name,disk.total,disk.used,disk.avail,disk.used_percent",
		"disk.used_percent:desc",
	)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	suffixes, err := listNodeNameSuffix()
	if err != nil {
		return err
	}
	fmt.Printf("valid data node suffixes: %s\n", strings.Join(suffixes, ", "))
	fmt.Printf("total data nodes: %d\n", len(suffixes))
	return nil
}

func catNodes(h, s string) ([]byte, error) {
	resp, err := client.Cat.Nodes(client.Cat.Nodes.WithV(true),
		client.Cat.Nodes.WithH(h),
		client.Cat.Nodes.WithS(s),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func listNodeFSDetails() error {
	b, err := getNodeStats("fs", "")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func getNodeStats(metric, filterPath string) ([]byte, error) {
	resp, err := client.Nodes.Stats(client.Nodes.Stats.WithHuman(),
		client.Nodes.Stats.WithPretty(),
		client.Nodes.Stats.WithMetric(metric),
		client.Nodes.Stats.WithFilterPath(filterPath),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func listNodesInfo() (*nodeInfoResp, error) {
	resp, err := client.Nodes.Info(client.Nodes.Info.WithPretty(),
		client.Nodes.Info.WithHuman(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	nodes := &nodeInfoResp{}
	if err := json.Unmarshal(b, nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func listNodesVersion() error {
	nodes, err := listNodesInfo()
	if err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "node\t elastic-version\t ip\t roles\t")
	for _, node := range nodes.Nodes {
		roles := strings.Join(node.Roles, " ")
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t\n", node.Name, node.Version, node.IP, roles)
	}
	w.Flush()
	return nil
}

func init() {
	listCmd.AddCommand(listNodesCmd)
	listNodesCmd.AddCommand(nodeStatsCmd, nodeSuffixesCmd, nodeStorageCmd, nodeFSDetailsCmd, nodeVersionCmd)
	// rootCmd.AddCommand(excludeNodeCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
