package authControllers

import (
	"ferdinand/app/models"
	authPages "ferdinand/views/pages/auth"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"gorm.io/gorm"
)

type SignInController struct {
	auth *caesarAuth.Auth
	db   *gorm.DB
}

func NewSignInController(auth *caesarAuth.Auth, db *gorm.DB) *SignInController {
	return &SignInController{
		auth: auth,
		db:   db,
	}
}

func (c *SignInController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.SignInPage())
}

type SignInValidator struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8"`
}

func (c *SignInController) Handle(ctx *caesar.Context) error {
	data, errors, ok := caesar.Validate[SignInValidator](ctx)
	if !ok {
		return ctx.Render(authPages.SignInForm(
			authPages.SignInInput{Email: data.Email, Password: data.Password},
			errors,
		))
	}

	var user models.User
	if err := c.db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Render(authPages.SignInForm(
				authPages.SignInInput{Email: data.Email, Password: data.Password},
				map[string]string{"Auth": "Invalid credentials."},
			))
		}
		return err
	}

	if !caesarAuth.CheckPasswordHash(data.Password, user.Password) {
		return ctx.Render(authPages.SignInForm(
			authPages.SignInInput{Email: data.Email, Password: data.Password},
			map[string]string{"Auth": "Invalid credentials."},
		))
	}

	if err := c.auth.Authenticate(ctx, user); err != nil {
		return err
	}

	return ctx.Redirect(AFTER_AUTH_REDIRECT_TO)
}
