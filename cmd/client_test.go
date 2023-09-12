package cmd

import (
	"os"
	"testing"
)

func TestInitEsClient(t *testing.T) {
	_, err := client.Info()
	if err != nil {
		t.Errorf("client.Info() = %v", err)
	}
}

func init() {
	homedir, _ := os.UserHomeDir()
	cfgFile = homedir + "/.esctl-dev.yaml"
	initConfig()
	initEsClient()
}
