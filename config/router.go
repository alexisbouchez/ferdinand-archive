package config

import (
	"ferdinand/app/controllers"
	apiControllers "ferdinand/app/controllers/api"
	authControllers "ferdinand/app/controllers/auth"
	"ferdinand/app/middleware"

	"ferdinand/views/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/vexillum"
)

func NewRouter(
	auth *caesarAuth.Auth,
	signUpController *authControllers.SignUpController,
	signInController *authControllers.SignInController,
	signOutController *authControllers.SignOutController,
	forgotPasswordController *authControllers.ForgotPasswordController,
	resetPasswordController *authControllers.ResetPasswordController,
	authGitHubOAuthController *authControllers.GitHubOAuthController,
	billingController *controllers.BillingController,
	settingsController *controllers.SettingsController,
	stripeController *controllers.StripeController,
	domainsController *controllers.DomainsController,
	mailApiKeysController *controllers.MailApiKeysController,
	emailsController *apiControllers.EmailsController,
	vexillum *vexillum.Vexillum,
) *caesar.Router {
	router := caesar.NewRouter()

	// Middleware
	router.Use(auth.SilentMiddleware)
	router.Use(middleware.ViewMiddleware(vexillum))

	// Marketing pages
	router.Render("/", pages.LandingPage())
	router.Render("/pricing", pages.PricingPage())

	// Auth routes
	router.Get("/auth/sign_up", signUpController.Show)
	router.Post("/auth/sign_up", signUpController.Handle)

	router.Get("/auth/sign_in", signInController.Show)
	router.Post("/auth/sign_in", signInController.Handle)

	router.Post("/auth/sign_out", signOutController.Handle).Use(auth.AuthMiddleware)

	// OAuth-related routes
	router.Get("/auth/github/redirect", authGitHubOAuthController.Redirect)
	router.Get("/auth/github/callback", authGitHubOAuthController.Callback)

	// Forgot password routes
	router.Get("/auth/forgot_password", forgotPasswordController.Show)
	router.Post("/auth/forgot_password", forgotPasswordController.Handle)

	// Reset password routes
	router.Get("/auth/reset_password/{jwt}", resetPasswordController.Show)
	router.Post("/auth/reset_password/{jwt}", resetPasswordController.Handle)

	// Overview page
	router.Render("/overview", pages.OverviewPage()).Use(auth.AuthMiddleware)

	// Settings-related routes
	router.Render("/settings", pages.SettingsPage()).Use(auth.AuthMiddleware)
	router.Patch("/settings", settingsController.Update).Use(auth.AuthMiddleware)
	router.Delete("/settings", settingsController.Delete).Use(auth.AuthMiddleware)

	// Domain-related routes
	router.
		Get("/domains", domainsController.Index).
		Use(auth.AuthMiddleware)
	router.
		Get("/domains/{id}", domainsController.Show).
		Use(auth.AuthMiddleware)
	router.Post("/domains", domainsController.Store).Use(auth.AuthMiddleware)
	router.
		Delete("/domains/{id}", domainsController.Delete).
		Use(auth.AuthMiddleware)
	router.Post("/domains/check_dns/{id}", domainsController.CheckDNS).Use(auth.AuthMiddleware)

	// API key-related routes
	router.
		Get("/api_keys", mailApiKeysController.Index).
		Use(auth.AuthMiddleware)
	router.
		Post("/api_keys", mailApiKeysController.Store).
		Use(auth.AuthMiddleware)
	router.
		Patch("/api_keys/{id}", mailApiKeysController.Update).
		Use(auth.AuthMiddleware)
	router.
		Delete("/api_keys/{id}", mailApiKeysController.Delete).
		Use(auth.AuthMiddleware)

	// Billing-related routes
	router.
		Get("/billing", billingController.Show).
		Use(auth.AuthMiddleware).
		Use(vexillum.EnsureFeatureEnabledMiddleware("billing"))
	router.
		Get("/billing/manage", billingController.Manage).
		Use(auth.AuthMiddleware).
		Use(vexillum.EnsureFeatureEnabledMiddleware("billing"))
	router.
		Get("/billing/payment_method", billingController.InitiatePaymentMethodChange).
		Use(auth.AuthMiddleware).
		Use(vexillum.EnsureFeatureEnabledMiddleware("billing"))

	// Webhooks routes
	router.Post("/webhooks/stripe", stripeController.HandleWebhook)

	// API-related routes
	router.Post("/api/v1/emails", emailsController.Send)

	return router
}
