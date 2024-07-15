package config

import (
	"github.com/caesar-rocks/core"
	"github.com/caesar-rocks/orm"
)

// EnvironmentVariables is a struct that holds all the environment variables that need to be validated.
// For full reference, see: https://github.com/go-playground/validator.
type EnvironmentVariables struct {
	// APP_KEY is the key used for encryption and decryption.
	APP_KEY string `validate:"required"`

	// HTTP_ADDR is the address to listen on for incoming requests.
	HTTP_ADDR string `validate:"required"`

	// APP_URL is the URL of the application.
	APP_URL string `validate:"required"`

	// DBMS is the database management system to use ("postgres", "mysql", "sqlite").
	DBMS orm.DBMS `validate:"oneof=postgres mysql sqlite"`

	// DSN is the data source name, which is a connection string for the database.
	DSN string `validate:"required"`

	// GITHUB_OAUTH_KEY is the key for the GitHub OAuth application.
	GITHUB_OAUTH_KEY string

	// GITHUB_OAUTH_SECRET is the secret for the GitHub OAuth application.
	GITHUB_OAUTH_SECRET string

	// GITHUB_OAUTH_CALLBACK_URL is the callback URL for the GitHub OAuth application.
	GITHUB_OAUTH_CALLBACK_URL string

	// STRIPE_SECRET_KEY is the key for the Stripe API.
	STRIPE_SECRET_KEY string

	// STRIPE_PUBLIC_KEY is the public key for the Stripe API.
	STRIPE_PUBLIC_KEY string

	// 	GITHUB_APP_ID is the ID for the GitHub App.
	GITHUB_APP_ID string

	// GITHUB_APP_PRIVATE_KEY is the private key for the GitHub App.
	GITHUB_APP_PRIVATE_KEY string

	// GITHUB_APP_WEBHOOK_SECRET is the secret for the GitHub App webhook.
	GITHUB_APP_WEBHOOK_SECRET string

	// GITHUB_APP_PRIVATE_KEY_PATH is the path to the private key for the GitHub App.
	GITHUB_APP_PRIVATE_KEY_PATH string

	// DB_HOST is the host for the database.
	DB_HOST string

	// SMTP_ADDR is the address for the SMTP server.
	SMTP_ADDR string

	// SMTP_DOMAIN is the domain for the SMTP server.
	SMTP_DOMAIN string

	// SMTP_ENABLE_TLS determines if TLS should be enabled for the SMTP server.
	SMTP_ENABLE_TLS string `validate:"oneof=true false"`

	// SMTP_ADMIN_EMAIL is the email address for the SMTP server administrator.
	SMTP_ADMIN_EMAIL string `validate:"email"`
}

func ProvideEnvironmentVariables() *EnvironmentVariables {
	return core.ValidateEnvironmentVariables[EnvironmentVariables]()
}
