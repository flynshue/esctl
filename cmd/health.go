package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "show cluster health",
	Long:  `esctl get health`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showHealth()
	},
}

func showHealth() error {
	resp, err := client.Cluster.Health(client.Cluster.Health.WithPretty())
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
	getCmd.AddCommand(healthCmd)
}
