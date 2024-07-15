package util

import (
	"context"
	"ferdinand/app/middleware"
)

func IsFeatureActive(ctx context.Context, feature string) bool {
	flags := ctx.Value(middleware.CTX_KEY_FLAGS).(map[string]bool)
	return flags[feature]
}
