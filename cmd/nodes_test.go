package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func TestListNodes(t *testing.T) {
	if err := listNodes(); err != nil {
		t.Error(err)
	}
}

func TestListNodesNames(t *testing.T) {
	nodeSuffixes, err := listNodeNameSuffix()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("valid data node suffixes: %s\n", strings.Join(nodeSuffixes, ", "))
}

func TestListNodeStorage(t *testing.T) {
	if err := listNodeStorage(); err != nil {
		t.Error(err)
	}
}

func TestListNodeFSDetails(t *testing.T) {
	if err := listNodeFSDetails(); err != nil {
		t.Error(err)
	}
}

func TestGetNodeVersion(t *testing.T) {
	if err := listNodesVersion(); err != nil {
		t.Error(err)
	}
}

func TestListShardCount(t *testing.T) {
	if err := listShardCount(); err != nil {
		t.Error(err)
	}
}
