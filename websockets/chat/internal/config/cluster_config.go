package config

type ClusterNodeConfig struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type ClusterFailoverConfig struct {
	Enabled bool `json:"enabled"`
	// Time in milliseconds between heartbeats
	Heartbeat int `json:"heartbeat"`
	// Number of failed heartbeats before a leader election is initiated.
	VoteAfter int `json:"vote_after"`
	// Number of failures before a node is considered dead
	NodeFailures int `json:"node_failures"`
}

type ClusterConfig struct {
	// List of all members of the cluster, including this member
	Nodes                 []ClusterNodeConfig   `json:"nodes"`
	MainName              string                `json:"main_name"` // Name of this cluster node
	ClusterFailOverConfig ClusterFailoverConfig `json:"cluster_fail_over_config"`
}
