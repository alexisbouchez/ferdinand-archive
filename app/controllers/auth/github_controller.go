package authControllers

import (
	"ferdinand/app/models"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"gorm.io/gorm"
)

type GitHubOAuthController struct {
	auth *auth.Auth
	db   *gorm.DB
}

func NewGitHubOAuthController(auth *auth.Auth, db *gorm.DB) *GitHubOAuthController {
	return &GitHubOAuthController{auth, db}
}

func (c *GitHubOAuthController) Redirect(ctx *caesar.Context) error {
	return c.auth.Social.Use("github").Redirect(ctx)
}

func (c *GitHubOAuthController) Callback(ctx *caesar.Context) error {
	oauthUser, err := c.auth.Social.Use("github").Callback(ctx)
	if err != nil {
		return err
	}

	// Find or create user
	user := models.User{GitHubUserID: oauthUser.UserID, Email: oauthUser.Email, FullName: oauthUser.Name}
	if err := c.db.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := c.db.Create(&user).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Authenticate user
	if err := c.auth.Authenticate(ctx, user); err != nil {
		return err
	}

	return ctx.Redirect(AFTER_AUTH_REDIRECT_TO)
}
