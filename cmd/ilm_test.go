package cmd

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

var (
	testIlmPolicyName = "test-timeseries"
	testIlmPolicy     = `{
  "policy": {
    "phases": {
      "hot": {                      
        "actions": {
          "rollover": {
            "max_size": "50GB",     
            "max_age": "30d"
          }
        }
      },
      "delete": {
        "min_age": "90d",           
        "actions": {
          "delete": {}              
        }
      }
    }
  }
}`
	testIndexPrefix   = "timeseries-7d"
	testIndexSuffix   = "000001"
	testIndexTemplate = `{
	"index_patterns": ["%s-*"],                 
	"template": {
	  "settings": {
		"number_of_shards": 1,
		"number_of_replicas": 1,
		"index.lifecycle.name": "%s",    
		"index.lifecycle.rollover_alias": "%s"  
	  }
	}
  }`
)

func TestIlm_BuildIlmIndexName(t *testing.T) {
	got := buildIlmIndexName("test-filebeat-7d-7.11.2", "000001")
	want := "%3Ctest-filebeat-7d-7.11.2-%7Bnow%2Fd%7D-000001%3E"
	if got != want {
		t.Errorf("got %s;\n wanted %s\n", got, want)
	}
}

func TestIlm_BootstrapIlmIdx(t *testing.T) {
	if err := setupIlm(); err != nil {
		t.Error(err)
	}
	if err := bootstrapIlmIdx(testIndexPrefix, testIndexSuffix); err != nil {
		t.Error(err)
	}
	if err := listIndexSettingsSummary(testIndexPrefix + "-*"); err != nil {
		t.Error(err)
	}
}

func setupIlm() error {
	// teardown if resources already exist
	index := fmt.Sprintf("%s-*", testIndexPrefix)
	client.Indices.Delete([]string{index})
	client.Indices.DeleteTemplate(testIndexPrefix)
	client.ILM.DeleteLifecycle(testIlmPolicyName)

	// create test ILM Policy
	buf := bytes.NewBufferString(testIlmPolicy)
	resp, err := client.ILM.PutLifecycle(testIlmPolicyName,
		client.ILM.PutLifecycle.WithBody(buf),
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

	// create test template
	buf = bytes.NewBufferString(fmt.Sprintf(testIndexTemplate, testIndexPrefix, testIlmPolicyName, testIndexPrefix))
	resp, err = client.Indices.PutIndexTemplate(testIndexPrefix, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func TestIlm_ListIlmPolicies(t *testing.T) {
	if err := setupIlm(); err != nil {
		t.Error(err)
	}
	if err := listIlmPolicies("*"); err != nil {
		t.Error(err)
	}
}

func TestIlm_GetIlmPolicy(t *testing.T) {
	if err := setupIlm(); err != nil {
		t.Error(err)
	}
	b, err := getIlmPolicy(testIlmPolicyName)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
