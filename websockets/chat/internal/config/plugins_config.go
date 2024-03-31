package config

type PluginConfig struct {
	Enabled bool `json:"enabled"`
	// Unique service name
	Name string `json:"name"`
	// Microseconds to wait before timeout
	Timeout int `json:"timeout"`
	// Filters for RPC calls: when to call vs when to skip the call
	Filters PluginRpcFilterConfig `json:"filters"`
	// What should the server do if plugin failed: HTTP error code
	FailureCode int `json:"failure_code"`
	// HTTP Error message to go with the code
	FailureMessage string `json:"failure_message"`
	// Address of plugin server of the form "tcp://localhost:123" or "unix://path_to_socket_file"
	ServiceAddr string `json:"service_addr"`
}

// PluginRpcFilterConfig filters for an individual RPC call. Filter strings are formatted as follows:
// <comma separated list of packet names> ;
// <comma separated list of topics or topic types> ;
// <actions (combination of C U D)>
//
// For instance:
// "acc,login;;CU" - grab packets {acc} or {login}; no filtering by topic, Create or Update action
// "pub,pres;me,p2p;"
type PluginRpcFilterConfig struct {
	// Filter by packet name, topic type [or exact name - not supported yet]. 2D: "pub,pres;p2p,me"
	FireHost *string `json:"fire_host"`
	// Filter by CUD, [exact user name - not supported yet]. 1D: "C"
	Account *string `json:"account"`
	// Filter by CUD, topic type[, exact name]: "p2p;CU"
	Topic *string `json:"topic"`
	// Filter by CUD, topic type[, exact topic name, exact user name]: "CU"
	Subscription *string `json:"subscription"`
	// Filter by C.D, topic type[, exact topic name, exact user name]: "grp;CD"
	Message *string `json:"message"`
	// Call Find service
	Find bool `json:"find"`
}
