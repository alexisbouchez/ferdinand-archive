package controllers

import (
	"ferdinand/app/models"
	"ferdinand/views/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/ui/toast"
	"gorm.io/gorm"
)

// DomainsController is the controller for managing mail domains
type DomainsController struct {
	db *gorm.DB
}

func NewDomainsController(db *gorm.DB) *DomainsController {
	return &DomainsController{db}
}

func (c *DomainsController) Index(ctx *caesar.Context) error {
	// Retrieve the user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Get the list of domains owned by the user
	var domains []models.Domain
	if err := c.db.Where("user_id = ?", user.ID).Find(&domains).Error; err != nil {
		return err
	}

	return ctx.Render(pages.ListDomainsPage(domains))
}

type StoreMailDomainValidator struct {
	Domain string `form:"domain" validate:"required"`
}

func (c *DomainsController) Store(ctx *caesar.Context) error {
	// Retrieve the user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Validate the input
	data, _, ok := caesar.Validate[StoreMailDomainValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	// Create the domain in the database
	domain := &models.Domain{Domain: data.Domain, UserID: user.ID}
	if err := c.db.Create(domain).Error; err != nil {
		return err
	}

	toast.Success(ctx, "Mail domain created successfully.")

	return ctx.Redirect("/domains/" + domain.ID)
}

func (c *DomainsController) Show(ctx *caesar.Context) error {
	// Retrieve the current user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Retrieve the domain from the bun database, where the domain matches the input
	// and the user ID matches the current user
	var domain models.Domain
	if err := c.db.Where("id = ? AND user_id = ?", ctx.PathValue("id"), user.ID).First(&domain).Error; err != nil {
		return err
	}

	return ctx.Render(pages.ShowDomainPage(domain))
}

func (c *DomainsController) Delete(ctx *caesar.Context) error {
	// Retrieve the current user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Delete the domain from the bun database, where the domain matches the input
	// and the user ID matches the current user
	domain := models.Domain{ID: ctx.PathValue("id"), UserID: user.ID}
	if err := c.db.Delete(&domain).Error; err != nil {
		return err
	}

	toast.Success(ctx, "Mail domain deleted successfully.")

	return ctx.Redirect("/domains")
}

func (c *DomainsController) CheckDNS(ctx *caesar.Context) error {
	// Retrieve the current user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Retrieve the domain from the bun database, where the domain matches the input
	domain := models.Domain{ID: ctx.PathValue("id"), UserID: user.ID}
	if err := c.db.First(&domain).Error; err != nil {
		return err
	}

	// Check the DNS records
	if err := domain.CheckDNS(); err != nil {
		return err
	}

	// Save the domain with the update dns status
	if err := c.db.Save(&domain).Error; err != nil {
		return err
	}

	return ctx.Redirect("/domains/" + domain.ID)
}
