package cmd

import "testing"

func TestShards_ListShards(t *testing.T) {
	testCases := []struct {
		name       string
		sort       string
		idxPattern []string
	}{
		{"ascending", "asc", []string{"*"}},
		{"descending", "desc", []string{"*"}},
		{"idxPattern", "desc", []string{".kibana*"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shardSort = tc.sort
			if err := listShards(tc.idxPattern); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestShards_ListShardsNodeBigger(t *testing.T) {
	if err := listShardsNodeBigger("es-data-03", "1kb", []string{"*"}); err != nil {
		t.Error(err)
	}
}

func TestShards_ListShardsForNode(t *testing.T) {
	testCases := []struct {
		name       string
		node       string
		idxPattern []string
	}{
		{"allShardsForNode", "es-data-03", []string{"*"}},
		{"indexPatternForNode", "es-data-03", []string{".kibana*"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shards, err := listShardsForNode(tc.node, tc.idxPattern)
			if err != nil {
				t.Error(err)
			}
			printShards(shards)
		})
	}
}

func TestShards_DisableShardAllocations(t *testing.T) {
	if err := disableShardAllocations(); err != nil {
		t.Error(err)
	}
}

func TestShards_GetShardAllocations(t *testing.T) {
	if err := getShardAllocations(); err != nil {
		t.Error(err)
	}
}

func TestShards_EnableShardAllocations(t *testing.T) {
	if err := enableShardAllocations(); err != nil {
		t.Error(err)
	}
}
