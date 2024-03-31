package config

type MediaConfig struct {
	// The name of the handler to use for file uploads.
	HandlerName string `json:"handler_name"`
	// Maximum allowed size of an uploaded file
	MaxFileUploadSize int64 `json:"max_file_upload_size"`
	// Garbage collection timeout
	GcPeriod int `json:"gc_period"`
	// Number of entries to delete in one pass
	GcBlockSize int            `json:"gc_block_size"`
	FS          *MediaConfigFS `json:"fs"`
}

type MediaConfigFS struct {
	FileUploadDirectory string   `json:"file_upload_directory"`
	ServeURL            string   `json:"serve_url"`
	CorsOrigins         []string `json:"cors_origins"`
}
