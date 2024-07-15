package authControllers

import (
	"bytes"
	"ferdinand/app/models"
	"ferdinand/views/mails"
	authPages "ferdinand/views/pages/auth"
	"os"
	"time"

	caesar "github.com/caesar-rocks/core"
	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt"
	gomail "gopkg.in/mail.v2"
	"gorm.io/gorm"
)

type ForgotPasswordController struct {
	db *gorm.DB
}

func NewForgotPasswordController(db *gorm.DB) *ForgotPasswordController {
	return &ForgotPasswordController{db}
}

func (c *ForgotPasswordController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.ForgotPasswordPage())
}

type ForgotPasswordValidator struct {
	Email string `form:"email" validate:"required,email"`
}

func (c *ForgotPasswordController) Handle(ctx *caesar.Context) error {
	data, _, ok := caesar.Validate[ForgotPasswordValidator](ctx)
	if !ok {
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	user := models.User{Email: data.Email}
	if err := c.db.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Render(authPages.ForgotPasswordSuccessAlert())
		}
		return err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("APP_KEY")))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	res := mails.ForgotPasswordMail(
		user.FullName,
		os.Getenv("APP_URL")+"/auth/reset_password/"+tokenString,
	)
	if err := res.Render(ctx.Context(), &buf); err != nil {
		log.Error("failed to render forgot password email", "err", err)
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer(os.Getenv("SMTP_DOMAIN"), 465, "ferdinand", os.Getenv("SMTP_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		log.Error("failed to send forgot password email", "err", err)
	}

	return ctx.Render(authPages.ForgotPasswordSuccessAlert())
}
