package config

type PathsConfig struct {
	// HTTP(S) address:port to listen on for websocket and long polling clients. Either a
	// numeric or a canonical name, e.g. ":80" or ":https". Could include a host name, e.g.
	// "localhost:80".
	// Could be blank: if TLS is not configured, will use ":80", otherwise ":443".
	Listen string `json:"listen"`
	// URL path for exposing runtime stats. Disabled if the path is blank.
	Expvar string `json:"expvar"`
	// Base URL path where the streaming and large file API calls are served,
	Api string `json:"api"`
	// // URL path for mounting the directory with static files .
	StaticMount string `json:"static_mount"`
	// Local path to static files. All files in this path are made accessible by HTTP.
	StaticData string `json:"static_data"`
}
