/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	data     string
	fileName string
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:     "console METHOD ENDPOINT",
	Aliases: []string{"esc"},
	Short:   "Send HTTP requests Elasticsearch REST API",
	PreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
	Long: `Interact with the REST APIs of Elasticsearch using http requests. This is useful for sending http requests to elasticsearch when we don't have commands built out for it yet.
esctl console GET /my-index-000001
esctl esc GET /my-index-000001
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) < 2 {
			return fmt.Errorf("must supply method and endpoint")
		}
		method := args[0]
		endpoint := args[1]
		switch {
		case data != "":
			return console(method, endpoint, []byte(data))
		case fileName != "":
			f, err := os.Open(fileName)
			if err != nil {
				return err
			}
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			return console(method, endpoint, b)
		default:
			return console(method, endpoint, nil)
		}
	},
}

func console(method, endpoint string, data []byte) error {
	b, err := esc.Do(method, endpoint, data)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func init() {
	rootCmd.AddCommand(consoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	consoleCmd.Flags().StringVarP(&data, "data", "d", "", "data body to be sent with http request")
	consoleCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file that contains data to be sent with request. --data takes precedence")
}
