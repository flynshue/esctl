package cmd

import "github.com/spf13/cobra"

var listRolesCmd = &cobra.Command{
	Use:     "roles",
	Aliases: []string{"role"},
	Short:   "list roles",
}

func init() {
	listCmd.AddCommand(listRolesCmd)
}
