package authControllers

import (
	"ferdinand/app/models"
	authPages "ferdinand/views/pages/auth"
	"os"
	"time"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type ResetPasswordController struct {
	db *gorm.DB
}

func NewResetPasswordController(db *gorm.DB) *ResetPasswordController {
	return &ResetPasswordController{db}
}

func (c *ResetPasswordController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.ResetPasswordPage())
}

type ResetPasswordValidator struct {
	Password        string `form:"password" validate:"required,min=8"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
}

func (c *ResetPasswordController) Handle(ctx *caesar.Context) error {
	// Fetch the JWT token from the URL parameter
	tokenString := ctx.PathValue("jwt")

	// Parse and validate the JWT token
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("APP_KEY")), nil
	})
	if err != nil || !validateTokenClaims(claims) {
		return err
	}

	// Validate the submitted form data
	data, errors, ok := caesar.Validate[ResetPasswordValidator](ctx)
	if !ok {
		return ctx.Render(authPages.ResetPasswordForm(errors))
	}

	// Fetch the user by ID from the token claims
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return err
	}

	var user *models.User
	if err := c.db.Where("id = ?", userID).First(user).Error; err != nil {
		return err
	}

	// Update the user's password in the database
	Password, err := auth.HashPassword(data.Password)
	if err != nil {
		return err
	}

	user.Password = Password
	if err := c.db.Save(user).Error; err != nil {
		return err
	}

	return ctx.Render(authPages.ResetPasswordSuccessAlert())
}

func validateTokenClaims(claims *jwt.MapClaims) bool {
	if exp, ok := (*claims)["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return false
		}
	}
	return true
}
