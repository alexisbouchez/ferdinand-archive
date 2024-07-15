package util

import (
	"context"
	"ferdinand/app/middleware"
)

func Route(ctx context.Context, suffix string) string {
	return "/orgs/" + ctx.Value(middleware.CTX_KEY_ORG_ID).(string) + suffix
}
