package config

type MediaConfig struct {
	// The name of the handler to use for file uploads.
	HandlerName string `json:"handler_name" mapstructure:"handler_name"`
	// Maximum allowed size of an uploaded file
	MaxFileUploadSize int64 `json:"max_file_upload_size" mapstructure:"max_file_upload_size"`
	// Garbage collection periodicity in seconds: unused or abandoned uploads are deleted.
	GcPeriod int `json:"gc_period" mapstructure:"gc_period"`
	// Number of entries to delete in one pass
	GcBlockSize int            `json:"gc_block_size" mapstructure:"gc_block_size"`
	FS          *MediaConfigFS `json:"fs" mapstructure:"fs"`
}

type MediaConfigFS struct {
	FileUploadDirectory string   `json:"file_upload_directory" mapstructure:"file_upload_directory"`
	CacheControl        string   `json:"cache_control" mapstructure:"cache_control"`
	ServerURL           string   `json:"server_url" mapstructure:"server_url"`
	CorsOrigins         []string `json:"cors_origins" mapstructure:"cors_origins"`
}
