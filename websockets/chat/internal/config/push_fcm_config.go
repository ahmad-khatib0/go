package config

type PushFcmConfig struct {
	Enabled         bool                     `json:"enabled"`
	DryRun          bool                     `json:"dry_run"`
	Credentials     PushFcmConfigCredentials `json:"credentials"`
	CredentialsFile string                   `json:"credentials_file"`
	TimeToLive      int                      `json:"time_to_live"`
	ApnsBundleID    string                   `json:"apns_bundle_id"` // Apple Push Notification service (APNs)
	Android         *PushCommonConfig        `json:"android"`
	Apns            *PushCommonConfig        `json:"apns"`
	WebPush         *PushCommonConfig        `json:"web_push"`
}

type PushFcmConfigCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}
