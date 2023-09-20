package cmd

import (
	"testing"
)

func TestCluster_GetClusterSettings(t *testing.T) {
	testCases := []struct {
		name       string
		filterPath string
	}{
		{"allSettings", ""},
		{"shardAllocations", "**.cluster.routing.allocation.enable"},
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

func TestCluster_SetClusterRebalance(t *testing.T) {
	if err := setRebalanceThrottle(300); err != nil {
		t.Error(err)
	}
}

func TestCluster_ResetRebalanceThrottle(t *testing.T) {
	if err := resetRebalanceThrottle(); err != nil {
		t.Error(err)
	}
}
