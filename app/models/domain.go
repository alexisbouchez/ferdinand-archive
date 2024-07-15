package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

const (
	DKIM_BITS_SIZE = 2048
)

type Domain struct {
	ID string `gorm:"primaryKey"`

	Domain string

	DKIMPrivateKey string
	DKIMPublicKey  string

	DNSVerified bool

	UserID string
	User   User

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExpectedDNSRecordType string

const (
	ExpectedDNSRecordTypeMX  ExpectedDNSRecordType = "MX"
	ExpectedDNSRecordTypeTXT ExpectedDNSRecordType = "TXT"
)

type ExpectedDNSRecord struct {
	Verified bool                  `json:"verified" example:"true" doc:"Record verification status"`
	Type     ExpectedDNSRecordType `json:"type" example:"MX" doc:"Record type"`
	Host     string                `json:"host" example:"example.com" doc:"Record host"`
	Value    string                `json:"value" example:"mail.example.com" doc:"Record value"`
}

func (d *Domain) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = xid.New().String()
	d.CreatedAt = time.Now()

	if err = d.assignDKIMKeysPair(); err != nil {
		return err
	}

	return
}

func (d *Domain) BeforeUpdate(tx *gorm.DB) (err error) {
	d.UpdatedAt = time.Now()
	return nil
}

func (d *Domain) assignDKIMKeysPair() error {
	key, err := rsa.GenerateKey(rand.Reader, DKIM_BITS_SIZE)
	if err != nil {
		return err
	}

	d.DKIMPrivateKey = exportRsaPrivateKeyAsStr(key)
	d.DKIMPublicKey, err = exportRsaPublicKeyAsStr(&key.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

func exportRsaPrivateKeyAsStr(privkey *rsa.PrivateKey) string {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	return base64.StdEncoding.EncodeToString(privkeyBytes)
}

func exportRsaPublicKeyAsStr(key *rsa.PublicKey) (string, error) {
	privkeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(privkeyBytes), nil
}

func (d *Domain) GetExpectedDNSRecords() []ExpectedDNSRecord {
	records := []ExpectedDNSRecord{
		{
			Type:  ExpectedDNSRecordTypeTXT,
			Host:  "ferdinand._domainkey." + d.Domain,
			Value: fmt.Sprintf("v=DKIM1; k=rsa; p=%s", d.DKIMPublicKey),
		},
		{
			Type:  ExpectedDNSRecordTypeTXT,
			Host:  "_dmarc." + d.Domain,
			Value: "v=DMARC1; p=none;",
		},
		{
			Type:  ExpectedDNSRecordTypeTXT,
			Host:  d.Domain,
			Value: fmt.Sprintf("v=spf1 mx a:%s -all", os.Getenv("SMTP_DOMAIN")),
		},
		{
			Type:  ExpectedDNSRecordTypeMX,
			Host:  d.Domain,
			Value: os.Getenv("SMTP_DOMAIN"),
		},
	}

	return records
}

func (d *Domain) CheckDNS() error {
	expectedRecords := d.GetExpectedDNSRecords()

	for i, record := range expectedRecords {
		log.Info("Checking DNS record", "type", record.Type, "host", record.Host, "value", record.Value)

		switch record.Type {
		case ExpectedDNSRecordTypeMX:
			log.Info("Looking up MX records", "domain", d.Domain)

			mxs, err := net.LookupMX(d.Domain)
			if err != nil {
				log.Error("Failed to lookup MX records", "err", err)
				expectedRecords[i].Verified = false
				continue
			}

			for _, mx := range mxs {
				log.Info("MX record", "host", mx.Host, "pref", mx.Pref, "value", record.Value)
				if strings.Contains(mx.Host, record.Value) {
					expectedRecords[i].Verified = true
					break
				}
			}
		case ExpectedDNSRecordTypeTXT:
			log.Info("Looking up TXT records", "domain", record.Host)

			txt, err := net.LookupTXT(record.Host)
			if err != nil {
				log.Error("Failed to lookup TXT records", "err", err)
				expectedRecords[i].Verified = false
				continue
			}

			for _, t := range txt {
				t = strings.Replace(t, "~all", "-all", 1)
				log.Info("TXT record", "value", t, "expected", record.Value, "verified", t == record.Value)
				if t == record.Value {
					expectedRecords[i].Verified = true
					break
				}
			}
		}
	}

	d.DNSVerified = true
	for _, record := range expectedRecords {
		if !record.Verified {
			d.DNSVerified = false
			break
		}
	}

	return nil
}
