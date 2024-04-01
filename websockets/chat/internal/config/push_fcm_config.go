package config

type PushFcmConfig struct {
	Enabled         bool                     `json:"enabled" mapstructure:"enabled"`
	DryRun          bool                     `json:"dry_run" mapstructure:"dry_run"`
	Credentials     PushFcmConfigCredentials `json:"credentials" mapstructure:"credentials"`
	CredentialsFile string                   `json:"credentials_file" mapstructure:"credentials_file"`
	TimeToLive      int                      `json:"time_to_live" mapstructure:"time_to_live"`
	// Apple Push Notification service (APNs)
	ApnsBundleID string           `json:"apns_bundle_id" mapstructure:"apns_bundle_id"`
	Android      PushCommonConfig `json:"android" mapstructure:"android"`
	Apns         PushCommonConfig `json:"apns" mapstructure:"apns"`
	WebPush      PushCommonConfig `json:"web_push" mapstructure:"web_push"`
}

type PushFcmConfigCredentials struct {
	Type                    string `json:"type" mapstructure:"type"`
	ProjectID               string `json:"project_id" mapstructure:"project_id"`
	PrivateKeyID            string `json:"private_key_id" mapstructure:"private_key_id"`
	PrivateKey              string `json:"private_key" mapstructure:"private_key"`
	ClientEmail             string `json:"client_email" mapstructure:"client_email"`
	AuthUri                 string `json:"auth_uri" mapstructure:"auth_uri"`
	TokenUri                string `json:"token_uri" mapstructure:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url" mapstructure:"auth_provider_x_509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url" mapstructure:"client_x_509_cert_url"`
}
