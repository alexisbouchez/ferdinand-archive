package config

import (
	"os"

	"github.com/caesar-rocks/vexillum"
)

func ProvideVexillum() *vexillum.Vexillum {
	v := vexillum.New()

	shouldActiveGitHubOAuth :=
		os.Getenv("GITHUB_OAUTH_KEY") != "" &&
			os.Getenv("GITHUB_OAUTH_SECRET") != "" &&
			os.Getenv("GITHUB_OAUTH_CALLBACK_URL") != ""
	if shouldActiveGitHubOAuth {
		v.Activate("github_oauth")
	}

	shouldActivateBilling :=
		os.Getenv("STRIPE_PUBLIC_KEY") != "" &&
			os.Getenv("STRIPE_SECRET_KEY") != ""
	if shouldActivateBilling {
		v.Activate("billing")
	}

	return v
}
