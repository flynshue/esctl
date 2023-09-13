package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "show cluster info",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

type clusterSettings struct {
	Cluster `json:"cluster"`
}

type Cluster struct {
	Routing `json:"routing"`
}

type Routing struct {
	Allocation `json:"allocation"`
}

type Allocation struct {
	Enable string `json:"enable"`
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

func init() {
	getCmd.AddCommand(clusterCmd)
}
