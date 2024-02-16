package cmd

type clusterSettings struct {
	Cluster `json:"cluster"`
	Indices `json:"indices,omitempty"`
	Action  `json:"action"`
}

type Action struct {
	Destructive string `json:"destructive_requires_name"`
}
type Cluster struct {
	Routing `json:"routing"`
}

type Indices struct {
	Recovery `json:"recovery,omitempty"`
}

type Recovery struct {
	MaxBps string `json:"max_bytes_per_sec,omitempty"`
}

type Routing struct {
	Allocation `json:"allocation,omitempty"`
}

type Allocation struct {
	Enable                      string `json:"enable,omitempty"`
	ClusterConcurrentRebalance  string `json:"cluster_concurrent_rebalance,omitempty"`
	NodeConcurrentIncomingRecov string `json:"node_concurrent_incoming_recoveries,omitempty"`
	NodeConcurrentOutgoingRecov string `json:"node_concurrent_outgoing_recoveries,omitempty"`
	NodeConcurrentRecov         string `json:"node_concurrent_recoveries,omitempty"`
	NodeInitialPriRecov         string `json:"node_initial_primaries_recoveries,omitempty"`
	Type                        string `json:"type,omitempty"`
	Exclude                     `json:"exclude,omitempty"`
}

type Exclude struct {
	Name string `json:"_name,omitempty"`
	IP   string `json:"_ip,omitempty"`
	Host string `json:"_host,omitempty"`
}

type WaterMarksResp struct {
	SingleDataNode              string `json:"cluster.routing.allocation.disk.watermark.enable_for_single_data_node"`
	FloodStage                  string `json:"cluster.routing.allocation.disk.watermark.flood_stage"`
	FloodStageFrozen            string `json:"cluster.routing.allocation.disk.watermark.flood_stage.frozen"`
	FloodStageFrozenMaxHeadRoom string `json:"cluster.routing.allocation.disk.watermark.flood_stage.frozen.max_headroom"`
	FloodStageMaxHeadRoom       string `json:"cluster.routing.allocation.disk.watermark.flood_stage.max_headroom"`
	High                        string `json:"cluster.routing.allocation.disk.watermark.high"`
	HighMaxHeadRoom             string `json:"cluster.routing.allocation.disk.watermark.high.max_headroom"`
	Low                         string `json:"cluster.routing.allocation.disk.watermark.low"`
	LowMaxHeadRoom              string `json:"cluster.routing.allocation.disk.watermark.low.max_headroom"`
}
