package config

import (
	errorPages "ferdinand/views/pages/errors"

	"github.com/caesar-rocks/core"
)

func ProvideErrorHandler() *core.ErrorHandler {
	return &core.ErrorHandler{Handle: func(ctx *core.Context, err error) {
		code := core.RetrieveErrorCode(err)

		if code == 404 {
			ctx.WithStatus(404).Render(errorPages.NotFoundPage())
		} else {
			ctx.WithStatus(code).Render(errorPages.ServerErrorPage(code))
		}
	}}
}
