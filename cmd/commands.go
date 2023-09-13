package cmd

import (
	"fmt"

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
	Use:   "get",
	Short: "get details for a resource",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list information for resource/s",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
}

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable resource/s",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
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
	rootCmd.AddCommand(disableCmd, getCmd, listCmd, enableCmd, versionCmd)
}
