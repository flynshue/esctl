package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

var (
	testRepo         = "test-elastic-fs"
	testRepoSettings = `{
	"type": "fs",
	"settings": {
	  "location": "/mnt/snapshot"
	}
  }`

	testSnapshotName = "test-snapshot-01"
	testSnapshotBody = `{
	  "indices": "%s",
	  "ignore_unavailable": true,
	  "include_global_state": false,
	  "metadata": {
		"taken_by": "user123",
		"taken_because": "backup before upgrading"
	  }
	}`
	testIdxs = []string{"test-idx-0001", "test-idx-0002"}
)

func TestSnapshots_ListSnapshotRepos(t *testing.T) {
	client.Snapshot.DeleteRepository([]string{testRepo})
	buf := bytes.NewBufferString(testRepoSettings)
	resp, err := client.Snapshot.CreateRepository(testRepo, buf)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
	if err := listSnapshotRepos(); err != nil {
		t.Error(err)
	}
}

func TestSnapshots_CatSnapshots(t *testing.T) {
	if err := snapshotSetup(); err != nil {
		t.Error(err)
	}
	if err := catSnapshots("*"); err != nil {
		t.Error(err)
	}
}

func snapshotSetup() error {
	// delete resources if they already exist
	client.Snapshot.Delete(testRepo, []string{testSnapshotName})
	client.Snapshot.DeleteRepository([]string{testRepo})
	client.Indices.Delete(testIdxs)

	for _, i := range testIdxs {
		client.Indices.Create(i)
	}
	buf := bytes.NewBufferString(testRepoSettings)
	_, err := client.Snapshot.CreateRepository(testRepo, buf)
	if err != nil {
		return err
	}
	buf = bytes.NewBufferString(fmt.Sprintf(testSnapshotBody, strings.Join(testIdxs, ",")))
	resp, err := client.Snapshot.Create(testRepo, testSnapshotName,
		client.Snapshot.Create.WithBody(buf),
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

func TestSnapshots_GetSnapshotRepos(t *testing.T) {
	client.Snapshot.DeleteRepository([]string{testRepo})
	buf := bytes.NewBufferString(testRepoSettings)
	_, err := client.Snapshot.CreateRepository(testRepo, buf)
	if err != nil {
		t.Error(err)
	}
	if err := getSnapshotRepos(testRepo); err != nil {
		t.Error(err)
	}
}

func TestSnapshots_GetSnapshot(t *testing.T) {
	if err := snapshotSetup(); err != nil {
		t.Error(err)
	}
	if err := getSnapshot(testSnapshotName); err != nil {
		t.Error(err)
	}
}

func TestSnapshots_ListSlmPolicies(t *testing.T) {
	if err := slmSetup(); err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 1)
	if err := listSlmPolicies(); err != nil {
		t.Error(err)
	}
}

func slmSetup() error {
	if err := snapshotSetup(); err != nil {
		return err
	}
	policyCfg := `{
		"schedule": "0 30 1 * * ?", 
		"name": "<%s-daily-snap-{now/d}>", 
		"repository": "%s", 
		"config": { 
		  "indices": "%s", 
		  "ignore_unavailable": false,
		  "include_global_state": false
		},
		"retention": { 
		  "expire_after": "30d", 
		  "min_count": 5, 
		  "max_count": 50 
		}
	  }`
	policyNames := []string{"test-slm-policy-01", "test-slm-policy-02"}
	for _, name := range policyNames {
		client.SlmDeleteLifecycle(name)
		buf := bytes.NewBufferString(fmt.Sprintf(policyCfg, name, testRepo, strings.Join(testIdxs, ",")))
		_, err := client.SlmPutLifecycle(name,
			client.SlmPutLifecycle.WithBody(buf),
		)
		if err != nil {
			return err
		}
		_, err = client.SlmExecuteLifecycle(name)
		if err != nil {
			return err
		}
	}
	return nil
}
