package templates

import (
	"context"
	"net/http"

	"rubensdev.com/go-image-processing/templates/manifest"
)

type contextKey string

var manifestManagerCtx = contextKey("manifestManager")
var isProdCtx = contextKey("isProd")

func WithManifestManager(mm *manifest.Manager, isProd bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), manifestManagerCtx, mm)
			ctx = context.WithValue(ctx, isProdCtx, isProd)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
