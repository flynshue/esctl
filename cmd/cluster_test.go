package cmd

import (
	"fmt"
	"testing"
)

func TestCluster_GetClusterSettings(t *testing.T) {
	testCases := []struct {
		name       string
		filterPath string
	}{
		{"allSettings", ""},
		{"shardAllocations", "**.cluster.routing.allocation.enable"},
		{"actionSettings", "**.action.destructive_requires_name"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := getClusterSettings(tc.filterPath)
			if err != nil {
				t.Error(err)
			}
			t.Log(string(b))
		})
	}
}

func TestCluster_GetClusterRebalance(t *testing.T) {
	if err := getClusterRebalance(); err != nil {
		t.Error(err)
	}
}

func TestCluster_GetDestructiveRequires(t *testing.T) {
	if err := getDestructiveRequires(); err != nil {
		t.Error(err)
	}
}

func TestCluster_SetClusterRebalance(t *testing.T) {
	if err := setRebalanceThrottle(300); err != nil {
		t.Error(err)
	}
}

func TestCluster_SetExcludeNode(t *testing.T) {
	testCases := []struct {
		name  string
		nodes string
	}{
		{"excludeSingleNode", "es-data-01"},
		{"excludeMultiNode", "es-data-01,es-data-02"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := setExcludeNode(tc.nodes); err != nil {
				t.Error(err)
			}
			fmt.Println()
			if err := getExcludedNodes(); err != nil {
				t.Error(err)
			}
			defer func() {
				setExcludeNode("")
			}()
		})
	}
}

func TestCluster_ClearExcludeNode(t *testing.T) {
	if err := setExcludeNode("es-data-02"); err != nil {
		t.Error(err)
	}
	getExcludedNodes()
	if err := setExcludeNode(""); err != nil {
		t.Error()
	}
	getExcludedNodes()
	listShardCount()
}

func TestCluster_ResetRebalanceThrottle(t *testing.T) {
	if err := resetRebalanceThrottle(); err != nil {
		t.Error(err)
	}
}

func TestCluster_ExplainClusterAllocation(t *testing.T) {
	if err := explainClusterAllocation(); err != nil {
		t.Error(err)
	}
}

func TestCluster_EnableDestructiveRequires(t *testing.T) {
	if err := enableDestructiveRequires(); err != nil {
		t.Error(err)
	}
}

func TestCluster_DisableDestructiveRequires(t *testing.T) {
	if err := disableDestructiveRequires(); err != nil {
		t.Error(err)
	}
}

func TestCluster_GetClusterWatermarks(t *testing.T) {
	if err := getClusterWatermarks(); err != nil {
		t.Error(err)
	}
}
