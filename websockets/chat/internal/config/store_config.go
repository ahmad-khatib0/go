package config

import "encoding/json"

type StoreConfig struct {
	// 16-byte key for XTEA. Used to initialize types.UidGenerator.
	UidKey []byte `json:"uid_key"`
	// Maximum number of results to return from adapter.
	MaxResults int `json:"max_results"`
	// DB adapter name to use.
	Adapters map[string]json.RawMessage `json:"adapters"`
}
