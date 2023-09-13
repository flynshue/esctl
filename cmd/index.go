/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var legacy bool

// indexCmd represents the index command
var getIndexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"idx"},
	Short:   "show information about one or more index",
}

// indexCmd represents the index command
var listIndexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"idx"},
	Short:   "show information about one or more index",
}

// idxSizesCmd represents the idxSizes command
var idxSizesCmd = &cobra.Command{
	Use:     "sizes",
	Aliases: []string{"size"},
	Short:   "show index sizes sorted (big -> small)",
	Long: `show index sizes sorted (big -> small)
ex: esctl list all indexes and their sizes
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showIdxSizes()
	},
}

// idxVersionCmd represents the idxVersion command
var idxVersionCmd = &cobra.Command{
	Use:   "versions [index pattern]",
	Short: "show index creation version",
	Long: `show index creation version
# list all indexes and their versions
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
	Use:     "template [name]",
	Aliases: []string{"templates"},
	Short:   "get one or more index templates",
	Long: `# List all index templates and their index patterns
esctl list index template

# Get list index templates that match template pattern
esctl list index template logs
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
	Long: `# Get details on for index template
esctl get index template .monitoring-beats
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply a template pattern")
		}
		return getIndexTemplate(args[0])
	},
}

type listIndexVersionResp struct {
	IndexSettings `json:"settings"`
}

type IndexSettings struct {
	Index `json:"index"`
}

type Index struct {
	IndexVersion `json:"version"`
}

type IndexVersion struct {
	Created string `json:"created_string"`
}

type indexTemplateLegacyResp map[string]indexTemplateSettings

type indexTemplateSettings struct {
	Patterns []string `json:"index_patterns"`
	Order    int      `json:"order"`
	Version  int      `json:"version"`
}

func showIdxSizes() error {
	resp, err := client.Cat.Indices(client.Cat.Indices.WithH("index,pri,rep,docs.count,store.size,pri.store.size"),
		client.Cat.Indices.WithHuman(),
		client.Cat.Indices.WithS("store.size:desc"),
		client.Cat.Indices.WithBytes("gb"),
		client.Cat.Indices.WithV(true),
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
	for k, idx := range idxs {
		fmt.Printf("index: %s, version:%s\n", k, idx.IndexVersion.Created)
	}
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
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
	listIndexCmd.AddCommand(idxSizesCmd, idxVersionCmd, listIndexTemplatesCmd)
	listIndexTemplatesCmd.Flags().BoolVar(&legacy, "legacy", false, "list only legacy index templates")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
