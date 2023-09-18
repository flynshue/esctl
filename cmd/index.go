/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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
	Use:   "versions [index pattern]",
	Short: "show index creation version",
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
			return listIndexDate("*")
		}
		return listIndexDate(args[0])
	},
}

func showIdxSizes() error {
	columns := "index,pri,rep,docs.count,store.size,pri.store.size"
	sort := "store.size:desc"
	b, err := catIndices(columns, sort, "*", "")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func listIndexDate(idxPattern string) error {
	columns := "index,pri,rep,docs.count,docs.deleted,store.size,creation.date"
	sort := "creation.date"
	b, err := catIndices(columns, sort, idxPattern, "json")
	if err != nil {
		return err
	}
	indices := []CatIndexResp{}
	err = json.Unmarshal(b, &indices)
	if err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "index\t primary_shards\t replica_shards\t docs\t deleted_docs\t store_size\t creation_date\t")
	for _, idx := range indices {
		date := parseCreateDate(idx.Date, localTime)
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t %s\t\n", idx.Index, idx.PrimaryShards, idx.ReplicaShards, idx.Docs, idx.DeletedDocs, idx.StoreSize, date)
	}
	w.Flush()
	return nil
}

func catIndices(columns, sort, idxPattern, format string) ([]byte, error) {
	resp, err := client.Cat.Indices(client.Cat.Indices.WithH(columns),
		client.Cat.Indices.WithS(sort),
		client.Cat.Indices.WithBytes("gb"),
		client.Cat.Indices.WithV(true),
		client.Cat.Indices.WithFormat(format),
		client.Cat.Indices.WithPretty(),
		client.Cat.Indices.WithIndex(idxPattern),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func listIndexVersion(pattern string) error {
	resp, err := client.Indices.Get([]string{pattern}, client.Indices.Get.WithHuman(),
		client.Indices.Get.WithHuman(),
		client.Indices.Get.WithExpandWildcards("all"),
		client.Indices.Get.WithFilterPath("*.settings.index.version.created_string"),
	)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("%s", resp.Status())
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	idxs := map[string]listIndexVersionResp{}
	if err := json.Unmarshal(b, &idxs); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "index\t version\t")
	for k, idx := range idxs {
		fmt.Fprintf(w, "%s\t %s\t\n", k, idx.IndexVersion.Created)
	}
	w.Flush()
	return nil
}

func listIndexTemplates(pattern string) error {
	resp, err := client.Cat.Templates(client.Cat.Templates.WithPretty(),
		client.Cat.Templates.WithName(pattern),
		client.Cat.Templates.WithV(true),
		client.Cat.Templates.WithS("name"),
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

func listIndexTemplatesLegacy(pattern string) error {
	resp, err := client.Indices.GetTemplate(client.Indices.GetTemplate.WithName(pattern),
		client.Indices.GetTemplate.WithHuman(),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var templates indexTemplateLegacyResp
	if err := json.Unmarshal(b, &templates); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "name\t index_pattern\t order\t")
	for name, t := range templates {
		fmt.Fprintf(w, "%s\t %v\t %d\t\n", name, t.Patterns, t.Order)
	}
	w.Flush()
	return nil
}

func getIndexTemplate(name string) error {
	resp, err := client.Indices.GetIndexTemplate(client.Indices.GetIndexTemplate.WithHuman(),
		client.Indices.GetIndexTemplate.WithPretty(),
		client.Indices.GetIndexTemplate.WithName(name),
	)
	if err != nil {
		return err
	}
	// legacy templates
	if resp.StatusCode == 404 {
		fmt.Fprintf(os.Stderr, "Warning: %s is a legacy index template. Legacy index templates have been deprecated starting in 7.8\n", name)
		resp, err = client.Indices.GetTemplate(client.Indices.GetTemplate.WithHuman(),
			client.Indices.GetTemplate.WithPretty(),
			client.Indices.GetTemplate.WithName(name),
		)
		if err != nil {
			return err
		}
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
	getCmd.AddCommand(getIndexCmd)
	getIndexCmd.AddCommand(getIndexTemplateCmd)
	listCmd.AddCommand(listIndexCmd)
	listIndexCmd.AddCommand(idxSizesCmd, idxVersionCmd, listIndexTemplatesCmd, listIndexDateCmd)
	listIndexTemplatesCmd.Flags().BoolVar(&legacy, "legacy", false, "list only legacy index templates")
	listIndexDateCmd.Flags().BoolVar(&localTime, "local", false, "display index creation timestamps in local time instead of UTC. Default is false.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
