package cmd

import "testing"

func TestListShards(t *testing.T) {
	testCases := []struct {
		name string
		sort string
	}{
		{"ascending", "asc"},
		{"descending", "desc"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shardSort = tc.sort
			if err := listShards(); err != nil {
				t.Error(err)
			}
		})
	}
}

// func TestListShardsJson(t *testing.T) {
// 	if err := listShardsJson(); err != nil {
// 		t.Error(err)
// 	}
// }

func TestListShardsForNode(t *testing.T) {
	if err := listShardsForNode("es-data-03"); err != nil {
		t.Error(err)
	}
}

func TestDisableShardAllocations(t *testing.T) {
	if err := disableShardAllocations(); err != nil {
		t.Error(err)
	}
}

func TestGetShardAllocations(t *testing.T) {
	if err := getShardAllocations(); err != nil {
		t.Error(err)
	}
}

func TestEnableShardAllocations(t *testing.T) {
	if err := enableShardAllocations(); err != nil {
		t.Error(err)
	}
}
