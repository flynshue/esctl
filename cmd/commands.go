package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable resource/s",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var getCmd = &cobra.Command{
	Use:   "get [command]",
	Short: "get details for a resource",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var listCmd = &cobra.Command{
	Use:   "list [command]",
	Short: "list information for resource/s",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var enableCmd = &cobra.Command{
	Use:   "enable [command]",
	Short: "enable resource/s",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var topCmd = &cobra.Command{
	Use:   "top [command]",
	Short: "Show elastic cluster stats",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var listCommands = &cobra.Command{
	Use:     "commands [command]",
	Short:   "List all the commands available",
	Aliases: []string{"cmd", "cmds", "command"},
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
		commandTree(w, cmd.Root())
		w.Flush()
	},
}

func commandTree(w *tabwriter.Writer, cmd *cobra.Command) {
	for _, c := range cmd.Commands() {
		if c.Name() == "help" {
			continue
		}
		fmt.Fprintf(w, "%s\t %s\t\n", c.UseLine(), c.Short)
		if c.HasSubCommands() {
			commandTree(w, c)
		}
	}
}

var (
	// GitCommit is updated with the Git tag by the Goreleaser build
	GitCommit = "unknown"
	// BuildDate is updated with the current ISO timestamp by the Goreleaser build
	BuildDate = "unknown"
	// Version is updated with the latest tag by the Goreleaser build
	Version = "unreleased"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the client version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t", Version)
		fmt.Println("Git commit:\t", GitCommit)
		fmt.Println("Date:\t\t", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(disableCmd, getCmd, listCmd, enableCmd, versionCmd, topCmd, listCommands)
}
