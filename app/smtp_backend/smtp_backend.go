package smtpBackend

import (
	"errors"
	mailBuilder "ferdinand/app/mail_builder"
	"ferdinand/app/models"
	smtpSender "ferdinand/app/smtp_sender"
	"ferdinand/util"
	"fmt"
	"io"
	"net/mail"
	"os"

	"github.com/charmbracelet/log"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"gorm.io/gorm"
)

// The Backend implements SMTP server methods.
type Backend struct {
	db *gorm.DB
}

// New creates a new Backend.
func New(db *gorm.DB) *Backend {
	return &Backend{db: db}
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	fmt.Println("New session", c.Hostname())
	return &Session{
		db: bkd.db,
	}, nil
}

// A Session is returned after successful login.
type Session struct {
	db *gorm.DB

	from   string
	to     string
	apiKey *models.APIKey
	domain *models.Domain
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is supported in this example.
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		// Check if the user exists
		if username != "ferdinand" {
			return errors.New("invalid username")
		}

		// Check if the API key is valid
		var apiKey models.APIKey
		err := s.db.Where("value = ?", password).First(&apiKey).Error
		if err != nil {
			return errors.New("invalid api key")
		}

		s.apiKey = &apiKey

		return nil
	}), nil
}

// Mail is called after MAIL FROM.
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// Check if the domain is verified, and if the user is allowed to send from it
	d, err := util.GetEmailDomain(from)
	if err != nil {
		log.Error("failed to get email domain", "error", err)
		return err
	}

	var domain models.Domain
	if err := s.db.Where(
		"domain = ? AND dns_verified = true AND user_id = ?",
		d, s.apiKey.UserID,
	).First(&domain).Error; err != nil {
		log.Error("failed to find domain", "domain", d, "error", err)
		return err
	}

	// Save the domain and the sender
	s.from = from
	s.domain = &domain

	return nil
}

// Rcpt is called after RCPT TO.
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// Save the recipient
	s.to = to

	return nil
}

// Data is called after DATA.
func (s *Session) Data(r io.Reader) error {
	// Read the message
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}

	// Build the email to send
	builder := mailBuilder.New(msg)
	outputMsg, err := builder.Build()
	if err != nil {
		return err
	}

	// Sign the email with DKIM
	outputMsg, err = builder.SignWithDKIM(outputMsg, s.domain.Domain, s.domain.DKIMPrivateKey)
	if err != nil {
		return err
	}

	// Send the email
	sender := smtpSender.New(os.Getenv("SMTP_DOMAIN"))
	sender.Send(s.from, s.to, outputMsg)

	return nil
}

// Reset is called after RSET.
func (s *Session) Reset() {}

// Logout is called after QUIT.
func (s *Session) Logout() error {
	return nil
}
