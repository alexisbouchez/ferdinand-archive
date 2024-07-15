package controllers

import (
	"ferdinand/app/models"

	"ferdinand/views/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/ui/toast"
	"gorm.io/gorm"
)

type SettingsController struct {
	db *gorm.DB
}

func NewSettingsController(db *gorm.DB) *SettingsController {
	return &SettingsController{db}
}

func (c *SettingsController) Edit(ctx *caesar.Context) error {
	return ctx.Render(pages.SettingsPage())
}

type SettingsValidator struct {
	Email           string `form:"email" validate:"required,email"`
	FullName        string `form:"full_name" validate:"required,min=3"`
	NewPassword     string `form:"new_password" validate:"omitempty,min=8"`
	ConfirmPassword string `form:"confirm_password" validate:"eqfield=NewPassword"`
}

func (c *SettingsController) Update(ctx *caesar.Context) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	data, errors, ok := caesar.Validate[SettingsValidator](ctx)
	if !ok {
		return ctx.Render(pages.SettingsForm(errors))
	}

	user.Email = data.Email
	user.FullName = data.FullName
	if data.NewPassword != "" {
		hashedPassword, err := caesarAuth.HashPassword(data.NewPassword)
		if err != nil {
			return caesar.NewError(400)
		}

		user.Password = hashedPassword
	}

	if err := c.db.Save(user).Error; err != nil {
		return err
	}

	toast.Success(ctx, "Settings updated successfully.")

	return ctx.Render(pages.SettingsForm(nil))
}

func (c *SettingsController) Delete(ctx *caesar.Context) error {
	// Retrieve the user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Delete the user
	if err := c.db.Delete(user).Error; err != nil {
		return err
	}

	// TODO: Bill the user, if they have a subscription (by emitting an event)

	return ctx.Redirect("/auth/sign_up")
}
