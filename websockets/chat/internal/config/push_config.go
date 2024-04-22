package config

type PushConfig struct {
	Name string
	FCM  *PushFcmConfig
}

// PushCommonConfigPayload  to be sent for a specific notification type.
type PushCommonConfigPayload struct {
	// Common for APNS and Android
	Body         string   `json:"body,omitempty" mapstructure:"body"`
	Title        string   `json:"title,omitempty" mapstructure:"title"`
	TitleLocKey  string   `json:"title_loc_key,omitempty" mapstructure:"title_loc_key"`
	TitleLocArgs []string `json:"title_loc_args,omitempty" mapstructure:"title_loc_args"`

	// Android
	BodyLocKey  string   `json:"body_loc_key,omitempty" mapstructure:"body_loc_key"`
	BodyLocArgs []string `json:"body_loc_args,omitempty" mapstructure:"body_loc_args"`
	Icon        string   `json:"icon,omitempty" mapstructure:"icon"`
	Color       string   `json:"color,omitempty" mapstructure:"color"`
	ClickAction string   `json:"click_action,omitempty" mapstructure:"click_action"`
	Sound       string   `json:"sound,omitempty" mapstructure:"sound"`
	Image       string   `json:"image,omitempty" mapstructure:"image"`

	// APNS
	Action          string   `json:"action,omitempty" mapstructure:"action"`
	ActionLocKey    string   `json:"action_loc_key,omitempty" mapstructure:"action_loc_key"`
	LaunchImage     string   `json:"launch_image,omitempty" mapstructure:"launch_image"`
	LocArgs         []string `json:"loc_args,omitempty" mapstructure:"loc_args"`
	LocKey          string   `json:"loc_key,omitempty" mapstructure:"loc_key"`
	Subtitle        string   `json:"subtitle,omitempty" mapstructure:"subtitle"`
	SummaryArg      string   `json:"summary_arg,omitempty" mapstructure:"summary_arg"`
	SummaryArgCount int      `json:"summary_arg_count,omitempty" mapstructure:"summary_arg_count"`
}

// PushCommonConfig is the configuration of a Notification payload.
type PushCommonConfig struct {
	Enabled bool `json:"enabled,omitempty" mapstructure:"enabled"`
	// Common defaults for all push types.
	PushCommonConfigPayload `mapstructure:"push_common_config_payload"`
	// Configs for specific push types.
	Msg PushCommonConfigPayload `json:"msg,omitempty" mapstructure:"msg"`
	Sub PushCommonConfigPayload `json:"sub,omitempty" mapstructure:"sub"`
}
