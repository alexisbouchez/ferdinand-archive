package controllers

import (
	"encoding/json"
	"ferdinand/app/models"
	"io"
	"net/http"
	"os"
	"time"

	caesar "github.com/caesar-rocks/core"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
	"gorm.io/gorm"
)

type StripeController struct {
	db *gorm.DB
}

func NewStripeController(db *gorm.DB) *StripeController {
	return &StripeController{db}
}

func (c *StripeController) HandleWebhook(ctx *caesar.Context) error {
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.ResponseWriter, ctx.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	event, err := webhook.ConstructEvent(
		payload,
		ctx.Request.Header.Get("Stripe-Signature"),
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
	)
	if err != nil {
		return err
	}

	switch event.Type {
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			return err
		}
		return c.handlePaymentMethodAttached(ctx, paymentMethod)
	case "payment_method.detached":
		stripeCustomerId := event.Data.PreviousAttributes["customer"].(string)
		return c.handlePaymentMethodDetached(ctx, stripeCustomerId)
	}

	return nil
}

func (c *StripeController) handlePaymentMethodAttached(ctx *caesar.Context, paymentMethod stripe.PaymentMethod) error {
	user := models.User{StripeCustomerID: paymentMethod.Customer.ID}
	if err := c.db.First(&user).Error; err != nil {
		return err
	}

	user.StripePaymentMethodID = paymentMethod.ID
	user.StripePaymentMethodExpirationDate = time.Date(
		int(paymentMethod.Card.ExpYear), time.Month(paymentMethod.Card.ExpMonth), 0, 0, 0, 0, 0, time.UTC,
	)
	if err := c.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (c *StripeController) handlePaymentMethodDetached(ctx *caesar.Context, stripeCustomerId string) error {
	// Retrieve the user by the Stripe customer ID
	user := models.User{StripeCustomerID: stripeCustomerId}
	if err := c.db.First(&user).Error; err != nil {
		return err
	}

	// Unset the payment method
	user.StripePaymentMethodID = ""
	user.StripePaymentMethodExpirationDate = time.Time{}
	if err := c.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
