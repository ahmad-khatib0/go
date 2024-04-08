package config

type PprofConf struct {
	// FileName to save profiling info
	FileName string `json:"file_name" mapstructure:"file_name"`
}
