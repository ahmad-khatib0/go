package config

type ClusterNodeConfig struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type ClusterFailoverConfig struct {
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// Time in milliseconds between heartbeats
	Heartbeat int `json:"heartbeat" mapstructure:"heartbeat"`
	// Number of failed heartbeats before a leader election is initiated.
	VoteAfter int `json:"vote_after" mapstructure:"vote_after"`
	// Number of failures before a node is considered dead
	NodeFailures int `json:"node_failures" mapstructure:"node_failures"`
}

type ClusterConfig struct {
	// List of all members of the cluster, including this member
	Nodes []ClusterNodeConfig `json:"nodes" mapstructure:"nodes"`
	// Name of this cluster node
	MainName        string                `json:"main_name" mapstructure:"main_name"`
	ClusterFailOver ClusterFailoverConfig `json:"cluster_fail_over" mapstructure:"cluster_fail_over"`
}
