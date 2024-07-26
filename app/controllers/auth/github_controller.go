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

	// Find user
	var user models.User
	if err := c.db.Where("github_user_id = ?", oauthUser.UserID).First(&user).Error; err != nil {
		// Authenticate user
		if err := c.auth.Authenticate(ctx, user); err != nil {
			return err
		}
	}

	// Or create user
	user = models.User{
		Email:        oauthUser.Email,
		FullName:     oauthUser.Name,
		GitHubUserID: oauthUser.UserID,
	}
	if err := c.db.Create(&user).Error; err != nil {
		return err
	}

	if err := c.auth.Authenticate(ctx, user); err != nil {
		return err
	}

	return ctx.Redirect(AFTER_AUTH_REDIRECT_TO)
}
