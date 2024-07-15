package config

import (
	"os"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

func ProvideStripe() *client.API {
	config := &stripe.BackendConfig{}
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), &stripe.Backends{
		API:     stripe.GetBackendWithConfig(stripe.APIBackend, config),
		Uploads: stripe.GetBackendWithConfig(stripe.UploadsBackend, config),
	})
	return sc
}
