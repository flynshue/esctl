/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	data     string
	fileName string
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:     "console [METHOD] [ENDPOINT]",
	Aliases: []string{"esc"},
	Short:   "Send HTTP requests Elasticsearch REST API",
	PreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
	Long: "Interact with the REST APIs of Elasticsearch using http requests. This is useful for sending http requests to elasticsearch when we don't have commands built out for it yet.",

	Example: `
# basic example
esctl console GET /my-index-000001

# command alias
esctl esc GET /my-index-000001

# without leading "/"
esctl esc GET my-index-000001

# supplying request data
esctl esc put /customer/_doc/1 -d \
'{
	"name": "John Doe"
}'

# supplying request data from file
esctl esc put /customer/_doc/2 -f /tmp/test-doc.json `,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) < 2 {
			return fmt.Errorf("must supply method and endpoint")
		}
		method := args[0]
		endpoint := args[1]
		switch {
		case data != "":
			return escConsole(method, endpoint, []byte(data))
		case fileName != "":
			b, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}
			return escConsole(method, endpoint, b)
		default:
			return escConsole(method, endpoint, nil)
		}
	},
}

// consoleCmd represents the console command
var kbnConsoleCmd = &cobra.Command{
	Use:   "kbn [METHOD] [ENDPOINT]",
	Short: "Send HTTP requests Kibana REST API",
	PreRun: func(cmd *cobra.Command, args []string) {
		initEsClient()
	},
	Long: "Interact with the REST APIs for Kibana using http requests.",

	Example: `
# basic example
esctl kbn GET /api/spaces/space
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) < 2 {
			return fmt.Errorf("must supply method and endpoint")
		}
		method := args[0]
		endpoint := args[1]
		switch {
		case data != "":
			return kbnConsole(method, endpoint, []byte(data))
		case fileName != "":
			b, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}
			return kbnConsole(method, endpoint, b)
		default:
			return kbnConsole(method, endpoint, nil)
		}
	},
}

func escConsole(method, endpoint string, data []byte) error {
	b, err := esc.Do(method, endpoint, data)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func kbnConsole(method, endpoint string, data []byte) error {
	b, err := kbn.Do(method, endpoint, data)
	if err != nil {
		return err
	}
	buf := bytes.Buffer{}
	if err := json.Indent(&buf, b, "", "  "); err != nil {
		return err
	}
	fmt.Println(buf.String())
	return nil
}

func init() {
	rootCmd.AddCommand(consoleCmd, kbnConsoleCmd)
	consoleCmd.Flags().StringVarP(&data, "data", "d", "", "data body to be sent with http request")
	consoleCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file that contains data to be sent with request. --data takes precedence")
	kbnConsoleCmd.Flags().StringVarP(&data, "data", "d", "", "data body to be sent with http request")
	kbnConsoleCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file that contains data to be sent with request. --data takes precedence")
}
