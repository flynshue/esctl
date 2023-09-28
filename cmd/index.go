/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func catIndices(columns, sort, format string, idxPattern []string) ([]byte, error) {
	if columns == "" {
		columns = "health,status,index,uuid,pri,rep,docs.count,docs.deleted,store.size,pri.store.size"
	}
	resp, err := client.Cat.Indices(client.Cat.Indices.WithH(columns),
		client.Cat.Indices.WithS(sort),
		client.Cat.Indices.WithV(true),
		client.Cat.Indices.WithFormat(format),
		client.Cat.Indices.WithPretty(),
		client.Cat.Indices.WithIndex(idxPattern...),
		client.Cat.Indices.WithHuman(),
		client.Cat.Indices.WithExpandWildcards("all"),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func deleteIndex(idxPattern []string) error {
	resp, err := client.Indices.Delete(idxPattern, client.Indices.Delete.WithHuman(),
		client.Indices.Delete.WithPretty(),
		client.Indices.Delete.WithExpandWildcards("open"),
		client.Indices.Delete.WithAllowNoIndices(true),
		client.Indices.Delete.WithIgnoreUnavailable(true),
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

func getIndexSettings(idxPattern string) ([]byte, error) {
	resp, err := client.Indices.GetSettings(client.Indices.GetSettings.WithIndex(idxPattern),
		client.Indices.GetSettings.WithExpandWildcards("all"),
		client.Indices.GetSettings.WithHuman(),
		client.Indices.GetSettings.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
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

func listIndexDate(idxPattern []string) error {
	columns := "index,pri,rep,docs.count,store.size,pri.store.size,creation.date"
	sort := "creation.date"
	b, err := catIndices(columns, sort, "json", idxPattern)
	if err != nil {
		return err
	}
	indices := []CatIndexResp{}
	err = json.Unmarshal(b, &indices)
	if err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "index\t pri\t rep\t docs.count\t pri.store.size\t store.size\t creation.date\t")
	for _, idx := range indices {
		date := parseUnixMilliDate(idx.Date, localTime)
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t %s\t\n", idx.Index, idx.PrimaryShards, idx.ReplicaShards, idx.Docs, idx.PriStoreSize, idx.StoreSize, date)
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

func listIndexSettingsSummary(idxPattern string) error {
	b, err := getIndexSettings(idxPattern)
	if err != nil {
		return err
	}
	idxs := make(map[string]listIndexSettingsResp)
	if err := json.Unmarshal(b, &idxs); err != nil {
		return err
	}
	w := newTabWriter()
	fmt.Fprintln(w, "index\t ilm_policy\t ilm_rollover_alias\t num_replicas\t num_shards\t auto_expand\t")
	for idx, s := range idxs {
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t\n", idx, s.Lifecycle.Name, s.Lifecycle.RolloverAlias, s.NumberOfReplicas, s.NumberOfShards, s.AutoExpand)
	}
	w.Flush()
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
	idxs := map[string]listIndexSettingsResp{}
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

func showIdxSizes(idxPattern []string) error {
	columns := "index,pri,rep,docs.count,pri.store.size,store.size"
	sort := "store.size:desc"
	b, err := catIndices(columns, sort, "", idxPattern)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func setIndexAutoExpand(index, autoExpand string) error {
	body := `{
		"index": {
			"auto_expand_replicas": "%s"
		}
	}`
	b, err := setIndexSettings(index, fmt.Sprintf(body, autoExpand))
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func setIndexReplicas(index string, rep int) error {
	body := `{
		"index": {
		  "number_of_replicas": %d

		 }
	   }`
	b, err := setIndexSettings(index, fmt.Sprintf(body, rep))
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func setIndexSettings(index, body string) ([]byte, error) {
	buf := bytes.NewBufferString(body)
	resp, err := client.Indices.PutSettings(buf, client.Indices.PutSettings.WithIndex(index))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
