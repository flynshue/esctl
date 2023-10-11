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
}
