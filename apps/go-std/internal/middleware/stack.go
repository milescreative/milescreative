package middleware

import (
	"net/http"

	"go-std/internal/config"
)

type MiddlewareContext struct {
	*config.App
}

func NewMiddlewareContext(app *config.App) *MiddlewareContext {
	return &MiddlewareContext{
		App: app,
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
