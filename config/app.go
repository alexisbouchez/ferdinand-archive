package config

import (
	"ferdinand/app/controllers"
	apiControllers "ferdinand/app/controllers/api"
	authControllers "ferdinand/app/controllers/auth"
	"ferdinand/public"
	"os"

	"github.com/caesar-rocks/core"
	"gorm.io/gorm"
)

func ProvideApp(db *gorm.DB) *core.App {
	app := core.NewApp(&core.AppConfig{
		Addr: os.Getenv("HTTP_ADDR"),
	})
	app.ErrorHandler = ProvideErrorHandler()

	auth := ProvideAuth(db)
	stripeClientAPI := ProvideStripe()
	vexillum := ProvideVexillum()

	router := NewRouter(
		auth,
		authControllers.NewSignUpController(auth, db),
		authControllers.NewSignInController(auth, db),
		authControllers.NewSignOutController(auth),
		authControllers.NewForgotPasswordController(db),
		authControllers.NewResetPasswordController(db),
		authControllers.NewGitHubOAuthController(auth, db),
		controllers.NewBillingController(stripeClientAPI),
		controllers.NewSettingsController(db),
		controllers.NewStripeController(db),
		controllers.NewDomainsController(db),
		controllers.NewMailApiKeysController(db),
		apiControllers.NewEmailsController(),
		vexillum,
	)

	core.ServeStaticAssets(public.FS)(router)
	app.Router = router

	return app
}
