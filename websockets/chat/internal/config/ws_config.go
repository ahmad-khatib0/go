package config

type WSConfig struct {
	// If true attempt to negotiate websocket per message compression (RFC 7692.4).
	// It should be disabled (set to false) if you are using MSFT IIS as a reverse proxy.
	WSCompressionEnabled bool `json:"ws_compression_enabled"`
	// Maximum message size allowed from client. Intended to prevent malicious client from sending
	// very large files inband (does not affect out of band uploads).
	MaxMessageSize int `json:"max_message_size"`
	// Maximum number of group topic subscribers.
	MaxSubscriberCount int `json:"max_subscriber_count"`
	// Masked tags: tags immutable on User (mask), mutable on Topic only within the mask.
	MaskedTagNamespaces []string `json:"masked_tag_namespaces"`
	// Maximum number of indexable tags.
	MaxTagCount int `json:"max_tag_count"`
}
