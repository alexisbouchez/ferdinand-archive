package main

import (
	"crypto/tls"
	smtpBackend "ferdinand/app/smtp_backend"
	"ferdinand/config"
	"os"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/caesar-rocks/core"
	"github.com/charmbracelet/log"
	"github.com/emersion/go-smtp"
)

func main() {
	// Create a new database connection
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Validate environment variables
	core.ValidateEnvironmentVariables[config.EnvironmentVariables]()

	// Create the SMTP backend
	backend := smtpBackend.New(db)

	srv := smtp.NewServer(backend)
	srv.Addr = os.Getenv("SMTP_ADDR")
	srv.Domain = os.Getenv("SMTP_DOMAIN")
	srv.WriteTimeout = 10 * time.Second
	srv.ReadTimeout = 10 * time.Second
	srv.MaxMessageBytes = 1024 * 1024
	srv.MaxRecipients = 50
	srv.AllowInsecureAuth = false

	// Set the TLS configuration
	if os.Getenv("SMTP_ENABLE_TLS") == "true" {
		certmagic.DefaultACME.Email = os.Getenv("SMTP_ADMIN_EMAIL")
		tlsConfig, err := certmagic.TLS([]string{os.Getenv("SMTP_DOMAIN")})
		if err != nil {
			log.Fatal("failed to get TLS configuration", "err", err)
		}
		tlsConfig.ClientAuth = tls.RequestClientCert
		tlsConfig.NextProtos = []string{"smtp", "smtps"}

		srv.TLSConfig = tlsConfig
	}

	// Start the server
	log.Info("starting smtp server", "addr", srv.Addr)
	if os.Getenv("SMTP_ENABLE_TLS") == "true" {
		log.Fatal(srv.ListenAndServeTLS())
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}
