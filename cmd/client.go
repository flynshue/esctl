package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/flynshue/esctl/pkg/esapi"
	"github.com/spf13/viper"
)

var (
	client *es.Client
	esc    *esapi.Client
)

func initEsClient() {
	var err error
	hosts := viper.GetStringSlice("hosts")
	if len(hosts) == 0 {
		log.Println("must supply hosts")
		os.Exit(1)
	}
	cfg := es.Config{
		Addresses: hosts,
		Username:  viper.GetString("username"),
		Password:  viper.GetString("password"),
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: viper.GetBool("insecure"),
			},
		},
	}
	client, err = es.NewClient(cfg)
	if err != nil {
		log.Printf("error creating client; %s\n", err)
		os.Exit(1)
	}
	esc = esapi.NewClient(hosts[0])
	esc.SetAuth(esapi.BasicAuth{Username: viper.GetString("username"), Password: viper.GetString("password")})
	esc.Headers = map[string]string{"Content-Type": "application/json"}
	if viper.GetBool("insecure") {
		esc.SkipTLS()
	}
}
