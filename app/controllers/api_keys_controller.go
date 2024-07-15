package controllers

import (
	"ferdinand/app/models"
	"ferdinand/views/pages"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"gorm.io/gorm"
)

type MailApiKeysController struct {
	db *gorm.DB
}

func NewMailApiKeysController(db *gorm.DB) *MailApiKeysController {
	return &MailApiKeysController{db}
}

func (c *MailApiKeysController) Index(ctx *caesar.Context) error {
	// Retrieve the user from the context
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Get the list of domains and API keys owned by the user
	var domains []models.Domain
	if err := c.db.Where("user_id = ?", user.ID).Find(&domains).Error; err != nil {
		return err
	}

	// Get the list of API keys owned by the user
	var apiKeys []models.APIKey
	if err := c.db.Where("user_id = ?", user.ID).Find(&apiKeys).Error; err != nil {
		return err
	}

	return ctx.Render(pages.APIKeysPage(domains, apiKeys))
}

type StoreMailApiKeyValidator struct {
	Name     string `form:"name"`
	DomainID string `form:"domain_id"`
}

func (c *MailApiKeysController) Store(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[StoreMailApiKeyValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	var onboarding bool

	if data.Name == "" {
		data.Name = "Onboarding"
		onboarding = true
	}

	apiKey := &models.APIKey{Name: data.Name, UserID: user.ID, DomainID: data.DomainID}
	if err := c.db.Create(apiKey).Error; err != nil {
		return err
	}

	if onboarding {
		return ctx.Render(pages.AddApiKeyOnboarding(apiKey.Value))
	}

	return ctx.Redirect("/api_keys")
}

func (c *MailApiKeysController) Update(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Get the API key by the user ID and the API key ID
	var apiKey *models.APIKey
	if err := c.db.Where("user_id = ?", user.ID).First(apiKey, ctx.PathValue("id")).Error; err != nil {
		return err
	}

	// Validate the submitted form data
	data, _, ok := caesar.Validate[StoreMailApiKeyValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	// Update the API key with the new information
	apiKey.Name = data.Name
	apiKey.DomainID = data.DomainID
	if err := c.db.Save(apiKey).Error; err != nil {
		return err
	}

	return ctx.Redirect("/api_keys")
}

func (c *MailApiKeysController) Delete(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Delete the API key, with the user ID and the API key ID.
	if err := c.db.
		Where("id = ? AND user_id = ?", ctx.PathValue("id"), user.ID).
		Delete(&models.APIKey{}).
		Error; err != nil {
		return err
	}

	return ctx.Redirect("/api_keys")
}
