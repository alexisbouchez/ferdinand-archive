package config

import (
	"context"
	"ferdinand/app/models"
	"os"
	"time"

	"github.com/caesar-rocks/auth"
	"gorm.io/gorm"
)

const (
	// AUTH_REDIRECT_TO is the default redirect path for the auth middleware
	AUTH_REDIRECT_TO = "/auth/sign_in"
)

func ProvideAuth(db *gorm.DB) *auth.Auth {
	return auth.NewAuth(&auth.AuthCfg{
		Key:    os.Getenv("APP_KEY"),
		MaxAge: time.Hour * 24 * 30,
		UserProvider: func(ctx context.Context, userID any) (any, error) {
			var user models.User
			err := db.WithContext(ctx).First(&user, "id = ?", userID).Error
			if err != nil {
				return nil, err
			}
			return &user, nil
		},
		RedirectTo: AUTH_REDIRECT_TO,
		SocialProviders: &map[string]auth.SocialAuthProvider{
			"github": {
				Key:         os.Getenv("GITHUB_OAUTH_KEY"),
				Secret:      os.Getenv("GITHUB_OAUTH_SECRET"),
				CallbackURL: os.Getenv("GITHUB_OAUTH_CALLBACK_URL"),
				Scopes:      []string{"user:email", "read:user"},
			},
		},
	})
}
