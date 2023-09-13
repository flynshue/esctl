package cmd

type nodeInfoResp struct {
	Nodes map[string]nodeSettings `json:"nodes"`
}

type nodeSettings struct {
	Name    string   `json:"name"`
	IP      string   `json:"ip"`
	Version string   `json:"version"`
	Roles   []string `json:"roles"`
}

type nodeStatsResp struct {
	Nodes map[string]nodeStats `json:"nodes"`
}

type nodeStats struct {
	Name       string `json:"name"`
	IndexStats `json:"indices"`
}

type IndexStats struct {
	ShardStats `json:"shard_stats"`
}

type ShardStats struct {
	Total int `json:"total_count"`
}
