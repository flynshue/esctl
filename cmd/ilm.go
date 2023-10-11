package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var createIlmIdxCmd = &cobra.Command{
	Use:     "ilm-index [index prefix] [index suffix]",
	Aliases: []string{"ilm-idx"},
	Short:   "bootstrap initial ilm index and designate it as the write index for the rollover alias specified",
	Long: `bootstrap initial ilm index and designate it as the write index for the rollover alias specified.
By default, the initial ilm index will be created as <index-prefix-{now/d}-index-suffix> and will use the index prefix as the rollover alias.
	`,
	Example: `# bootstrap initial ilm index with name test-filebeat-7d-7.11.2-2023.03.27-000001
esctl create ilm-index test-filebeat-7d-7.11.2 000001	
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bootstrapIlmIdx(args[0], args[1])
	},
}

var listIlmPoliciesCmd = &cobra.Command{
	Use:   "ilm",
	Short: "list ilm policies",
	Example: `# list all ilm policies
esctl list ilm

# list ilm policies by policy name pattern
esctl list ilmtest-*
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listIlmPolicies("*")
		}
		return listIlmPolicies(args[0])
	},
}

var getIlmPolicyCmd = &cobra.Command{
	Use:   "ilm",
	Short: "get ilm policy details",
	Example: `# get ilm policy details specific policy id
esctl get ilm metrics

# get ilm policy details for policy pattern
esctl get ilm synth*
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must supply policy name or policy pattern")
		}
		b, err := getIlmPolicy(args[0])
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

func bootstrapIlmIdx(idxPrefix, idxSuffix string) error {
	reqBody := `{
		"aliases": {
		  "%s": {
			"is_write_index": true
		  }
		}
	  }`
	buf := bytes.NewBufferString(fmt.Sprintf(reqBody, idxPrefix))
	index := buildIlmIndexName(idxPrefix, idxSuffix)
	resp, err := client.Indices.Create(index,
		client.Indices.Create.WithBody(buf),
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

func buildIlmIndexName(idxPrefix, idxSuffix string) string {
	return fmt.Sprintf("%%3C%s-%%7Bnow%%2Fd%%7D-%s%%3E", idxPrefix, idxSuffix)
}

func listIlmPolicies(policy string) error {
	b, err := getIlmPolicy(policy)
	if err != nil {
		return err
	}
	ilmPolices := map[string]any{}
	if err := json.Unmarshal(b, &ilmPolices); err != nil {
		return err
	}
	fmt.Println("ilm_policy")
	for policy := range ilmPolices {
		fmt.Println(policy)
	}
	return nil
}

func getIlmPolicy(policy string) ([]byte, error) {
	resp, err := client.ILM.GetLifecycle(client.ILM.GetLifecycle.WithPretty(),
		client.ILM.GetLifecycle.WithHuman(),
		client.ILM.GetLifecycle.WithPolicy(policy),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func init() {
	createCmd.AddCommand(createIlmIdxCmd)
	getCmd.AddCommand(getIlmPolicyCmd)
	listCmd.AddCommand(listIlmPoliciesCmd)
}
