package models

import (
	"ferdinand/util"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type APIKey struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Value string

	// Optional relationship with Domain
	Domain   Domain `gorm:"foreignKey:DomainID"`
	DomainID string `gorm:"default:null"`

	User   User `gorm:"foreignKey:UserID"`
	UserID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (domain *APIKey) BeforeCreate(tx *gorm.DB) (err error) {
	domain.ID = xid.New().String()
	domain.CreatedAt = time.Now()
	value, err := util.GenerateSecretKey()
	if err != nil {
		return err
	}
	domain.Value = value
	return nil
}

func (domain *APIKey) BeforeUpdate(tx *gorm.DB) (err error) {
	domain.UpdatedAt = time.Now()

	return nil
}
