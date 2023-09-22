package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	legacy    bool
	localTime bool
)

// indexCmd represents the index command
var getIndexCmd = &cobra.Command{
	Use:     "index [command]",
	Aliases: []string{"idx"},
	Short:   "get detailed information about one or more index",
}

// indexCmd represents the index command
var listIndexCmd = &cobra.Command{
	Use:     "index [command]",
	Aliases: []string{"idx"},
	Short:   "show information about one or more index",
}

var setIndexCmd = &cobra.Command{
	Use:     "index [command]",
	Aliases: []string{"idx"},
	Short:   "set configuration on index",
}

// idxSizesCmd represents the idxSizes command
var idxSizesCmd = &cobra.Command{
	Use:     "sizes",
	Aliases: []string{"size"},
	Short:   "show index sizes sorted (big -> small)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showIdxSizes()
	},
}

// idxVersionCmd represents the idxVersion command
var idxVersionCmd = &cobra.Command{
	Use:     "versions [index pattern]",
	Aliases: []string{"version"},
	Short:   "show index creation version",
	Example: `# list all indexes and their versions
esctl list index versions

# list all indexes and their versions for pattern
esctl list index versions watch*
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := "*"
		if len(args) != 0 {
			pattern = args[0]
		}
		return listIndexVersion(pattern)
	},
}

var listIndexTemplatesCmd = &cobra.Command{
	Use:     "template [template name pattern]",
	Aliases: []string{"templates"},
	Short:   "get one or more index templates",
	Example: `# List all index templates and their index patterns
esctl list index template

# Get list index templates that match template pattern
esctl list index template .monit*
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := "*"
		if len(args) != 0 {
			pattern = args[0]
		}
		if legacy {
			return listIndexTemplatesLegacy("*")
		}
		return listIndexTemplates(pattern)
	},
}

var getIndexTemplateCmd = &cobra.Command{
	Use:     "template [name]",
	Aliases: []string{"templates"},
	Short:   "get details for index template",
	Example: `# Get details on for index template
esctl get index template .monitoring-beats
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply a template pattern")
		}
		return getIndexTemplate(args[0])
	},
}

var listIndexDateCmd = &cobra.Command{
	Use:   "date [idx Pattern]",
	Short: "list all indexes with their creation date",
	Example: `# List indexes and their creation date that match index pattern .fleet*
esctl list index date .fleet*
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listIndexDate([]string{"*"})
		}
		return listIndexDate(args)
	},
}

var listIndexSettingsCmd = &cobra.Command{
	Use:     "settings [index pattern]",
	Aliases: []string{"config", "cfg"},
	Short:   "list indexes with a summary of settings. Includes replicas, shards, ilm policy, ilm rollover alias, and auto expand replicas",
	Example: `# List all indexes with summary of settings
esctl list index settings

# List indexes matching pattern with summary of settings
esctl list index settings .fleet-*
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listIndexSettingsSummary("*")
		}
		return listIndexSettingsSummary(args[0])
	},
}

var getIndexSettingsCmd = &cobra.Command{
	Use:     "settings [index pattern]",
	Aliases: []string{"config", "cfg"},
	Short:   "get full details of settings for index/index pattern",
	Example: `# Get index settings details for specific index
esctl get index settings .fleet-file-data-agent-000001

# Get index settings details for index pattern
esctl get index settings .fleet-*
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply index or index pattern")
		}
		b, err := getIndexSettings(args[0])
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

var setIndexReplicasCmd = &cobra.Command{
	Use:     "replicas [index] [number of replicas]",
	Aliases: []string{"replica", "rep"},
	Short:   "set the number of replicas for an index",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply index and number of replicas")
		}
		idx := args[0]
		rep, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		return setIndexReplicas(idx, rep)
	},
}

var setIndexAutoExpandCmd = &cobra.Command{
	Use:   "auto-expand [index] [replica range|false]",
	Short: "Auto-expand the number of replicas based on the number of data nodes in the cluster.",
	Long: `Auto-expand the number of replicas based on the number of data nodes in the cluster.
Replica range is dash delimited: 0-1, default value is false.
	`,
	Example: `# Set auto-expand to 0-1 replicas.
esctl set auto-expand test-logs-0001 0-1

# Disable auto-expand replicas.  Useful if you need to manually set the replicas to 0.
esctl set auto-expand test-logs-0001 false
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("must supply index and auto expand value")
		}
		return setIndexAutoExpand(args[0], args[1])
	},
}

func init() {
	getCmd.AddCommand(getIndexCmd)
	setCmd.AddCommand(setIndexCmd)
	setIndexCmd.AddCommand(setIndexReplicasCmd, setIndexAutoExpandCmd)
	getIndexCmd.AddCommand(getIndexTemplateCmd, getIndexSettingsCmd)
	listCmd.AddCommand(listIndexCmd)
	listIndexCmd.AddCommand(idxSizesCmd, idxVersionCmd, listIndexTemplatesCmd, listIndexDateCmd, listIndexSettingsCmd)
	listIndexTemplatesCmd.Flags().BoolVar(&legacy, "legacy", false, "list only legacy index templates")
	listIndexDateCmd.Flags().BoolVar(&localTime, "local", false, "display index creation timestamps in local time instead of UTC. Default is false.")

}
