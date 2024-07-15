package util

import (
	"context"
	"ferdinand/app/middleware"
)

func GetClassForTab(tabPath string, currentPath string) string {
	if tabPath == currentPath {
		return "text-ferdinand-200 whitespace-nowrap"
	}

	return "hover:text-ferdinand-200 transition-colors whitespace-nowrap"
}

func RetrievePath(ctx context.Context) string {
	ctxValue := ctx.Value(middleware.CTX_KEY_PATH)
	if ctxValue == nil {
		return ""
	}
	path, ok := ctxValue.(string)
	if !ok {
		return ""
	}

	return path
}
