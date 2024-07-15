package models

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type User struct {
	ID string `gorm:"primaryKey"`

	Email    string `gorm:"unique;not null"`
	FullName string `gorm:"not null"`

	Password string

	GitHubUserID string `gorm:"unique"`

	StripeCustomerID                  string `gorm:"unique"`
	StripePaymentMethodID             string `gorm:"unique"`
	StripePaymentMethodExpirationDate time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = xid.New().String()
	user.CreatedAt = time.Now()
	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return
}

func (u *User) HasActivePaymentMethod() bool {
	if u.StripePaymentMethodID == "" {
		return false
	}

	return u.StripePaymentMethodExpirationDate.After(time.Now())
}
