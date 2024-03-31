package config

type PushConfig struct {
	Name string
	FCM  *PushFcmConfig
}

// PushCommonConfigPayload  to be sent for a specific notification type.
type PushCommonConfigPayload struct {
	// Common for APNS and Android
	Body         string   `json:"body,omitempty"`
	Title        string   `json:"title,omitempty"`
	TitleLocKey  string   `json:"title_loc_key,omitempty"`
	TitleLocArgs []string `json:"title_loc_args,omitempty"`

	// Android
	BodyLocKey  string   `json:"body_loc_key,omitempty"`
	BodyLocArgs []string `json:"body_loc_args,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Color       string   `json:"color,omitempty"`
	ClickAction string   `json:"click_action,omitempty"`
	Sound       string   `json:"sound,omitempty"`
	Image       string   `json:"image,omitempty"`

	// APNS
	Action          string   `json:"action,omitempty"`
	ActionLocKey    string   `json:"action_loc_key,omitempty"`
	LaunchImage     string   `json:"launch_image,omitempty"`
	LocArgs         []string `json:"loc_args,omitempty"`
	LocKey          string   `json:"loc_key,omitempty"`
	Subtitle        string   `json:"subtitle,omitempty"`
	SummaryArg      string   `json:"summary_arg,omitempty"`
	SummaryArgCount int      `json:"summary_arg_count,omitempty"`
}

// PushCommonConfig is the configuration of a Notification payload.
type PushCommonConfig struct {
	Enabled bool `json:"enabled,omitempty"`
	// Common defaults for all push types.
	PushCommonConfigPayload
	// Configs for specific push types.
	Msg PushCommonConfigPayload `json:"msg,omitempty"`
	Sub PushCommonConfigPayload `json:"sub,omitempty"`
}
