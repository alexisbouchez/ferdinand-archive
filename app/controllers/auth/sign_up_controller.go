package authControllers

import (
	"errors"
	"ferdinand/app/models"
	authPages "ferdinand/views/pages/auth"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type SignUpController struct {
	auth *caesarAuth.Auth
	db   *gorm.DB
}

func NewSignUpController(auth *caesarAuth.Auth, db *gorm.DB) *SignUpController {
	return &SignUpController{auth, db}
}

func (c *SignUpController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.SignUpPage())
}

type SignUpValidator struct {
	Email    string `form:"email" validate:"required,email"`
	FullName string `form:"full_name" validate:"required,min=3"`
	Password string `form:"password" validate:"required,min=8"`
}

func (c *SignUpController) Handle(ctx *caesar.Context) error {
	data, validationErrors, ok := caesar.Validate[SignUpValidator](ctx)
	if !ok {
		return ctx.Render(authPages.SignUpForm(
			authPages.SignUpInput{Email: data.Email, FullName: data.FullName, Password: data.Password},
			validationErrors,
		))
	}

	user := models.User{Email: data.Email, FullName: data.FullName, Password: data.Password}
	if err := c.db.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return ctx.Render(authPages.SignUpForm(
				authPages.SignUpInput{Email: data.Email, FullName: data.FullName, Password: data.Password},
				map[string]string{"email": "Email is already in use"},
			))
		}
		log.Error("error while inserting user into the database", "err", err)
	}

	if err := c.auth.Authenticate(ctx, user); err != nil {
		log.Info("error authenticating user", "err", err)
		return caesar.NewError(400)
	}

	return ctx.Redirect(AFTER_AUTH_REDIRECT_TO)
}
