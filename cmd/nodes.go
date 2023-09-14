/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var listNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "show information about one or more node",
}

var nodeStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "list ES nodes with usage statistics",
	Long: `list ES nodes with usage statistics
esctl get nodes stats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodes()
	},
}

var nodeSuffixesCmd = &cobra.Command{
	Use:   "suffix",
	Short: "list ES nodes name suffixes",
	Long: `list ES nodes name suffixes
esctl get nodes suffix`,
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
	Long: `list ES nodes HDD usage
esctl get nodes storage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodeStorage()
	},
}

var nodeFSDetailsCmd = &cobra.Command{
	Use:     "filesystem",
	Aliases: []string{"fs"},
	Short:   "list ES nodes filesystem details",
	Long: `list ES nodes HDD usage
esctl get nodes filesystem
esctl get nodes fs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNodeFSDetails()
	},
}

var nodeVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "list ES nodes version",
	Long:  "esctl get nodes version",
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "node\t elastic-version\t ip\t roles\t")
	for _, node := range nodes.Nodes {
		roles := strings.Join(node.Roles, "")
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t\n", node.Name, node.IP, node.Version, roles)
	}
	w.Flush()
	return nil
}

func init() {
	listCmd.AddCommand(listNodesCmd)
	listNodesCmd.AddCommand(nodeStatsCmd, nodeSuffixesCmd, nodeStorageCmd, nodeFSDetailsCmd, nodeVersionCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
