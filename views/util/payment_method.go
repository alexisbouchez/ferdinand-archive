package util

import (
	"context"
	"ferdinand/app/middleware"
)

func ShowPaymentMethodDialog(ctx context.Context) bool {
	return ctx.Value(middleware.CTX_KEY_SHOW_PAYMENT_METHOD_DIALOG).(bool)
}
