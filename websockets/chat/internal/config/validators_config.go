package config

type ValidatorConfig struct {
	Email *ValidatorConfigEmail `json:"email" mapstructure:"email"`
}

type ValidatorConfigEmail struct {
	AddToTags bool `json:"add_to_tags" mapstructure:"add_to_tags"`
	//  Authentication level which triggers this validator: "auth", "anon"... or ""
	Required []string `json:"required" mapstructure:"required"`
	// Address of the host where the chat server is running. This will be used in URLs in the email.
	HostUrl string `json:"host_url" mapstructure:"host_url"`
	// Address of the SMPT server to use.
	SmtpServer string `json:"smtp_server" mapstructure:"smtp_server"`
	// SMTP port to use. "25" for basic email RFC 5321 (2821, 821), "587" for RFC 3207 (TLS).
	SmtpPort int `json:"smtp_port" mapstructure:"smtp_port"`
	// RFC 5322 email address to show in the From: field.
	Sender string `json:"sender" mapstructure:"sender"`
	// Optional login to use for authentication; if missing, the connection is not authenticated.
	Login string `json:"login" mapstructure:"login"`
	// Password to use when authenticating the sender; used only if "login" is provided
	Password string `json:"password" mapstructure:"password"`
	// Authentication mechanism to use, optional. One of "login", "cram-md5", "plain" (default).
	AuthMechanism string `json:"auth_mechanism" mapstructure:"auth_mechanism"`
	// FQDN to use in SMTP HELO/EHLO command; if missing, the hostname from "host_url" is used.
	SmtpHeloHost string `json:"smtp_helo_host" mapstructure:"smtp_helo_host"`
	// Skip verification of the server's certificate chain and host name.
	// In this mode, TLS is susceptible to machine-in-the-middle attacks.
	InsecureSkipVerify bool `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify"`
	// Allow this many confirmation attempts before blocking the credential.
	MaxRetries int `json:"max_retries" mapstructure:"max_retries"`
	// List of email domains allowed to be used for registration.
	// Missing or empty list means any email domain is accepted.
	Domains []string `json:"domains" mapstructure:"domains"`
	// Dummy response to accept.
	//
	// === IMPORTANT ===
	//
	// REMOVE IN PRODUCTION!!! Otherwise anyone will be able to register
	// with fake emails.
	DebugResponse string `json:"debug_response" mapstructure:"debug_response"`
}
