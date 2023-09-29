package cmd

import (
	"bytes"
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

func init() {
	createCmd.AddCommand(createIlmIdxCmd)
}
