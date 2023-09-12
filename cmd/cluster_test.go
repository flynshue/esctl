package cmd

import (
	"testing"
)

func TestGetClusterSettings(t *testing.T) {
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
