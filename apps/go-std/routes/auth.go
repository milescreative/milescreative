package routes

import (
	"go-std/internal/auth"
	"go-std/internal/config"
	"net/http"
)

func AuthRoutes(mux *http.ServeMux, app *config.App) {

	router := NewRouter_(mux)
	r := router

	a := auth.NewAuthHandlers(app)

	//todo: move to root
	r.GET("/{$}", a.EnvHandler)

	r.GET("/api/auth/login/google", a.LoginHandler)
	r.GET("/api/auth/callback/google", a.CallbackHandler)
	r.GET("/api/auth/validate", a.ValidateSessionHandler)
	r.POST("/api/auth/refresh", a.RefreshTokenHandler)
	r.POST("/api/auth/logout", a.LogoutHandler)
	r.GET("/api/auth/csrf", a.GetCSRFTokenHandler)
	r.GET("/api/auth/user", a.GetUserHandler)
	r.POST("/api/auth/user", a.UpdateUserHandler)
	r.DELETE("/api/auth/user", a.DeleteUserHandler)
	r.GET("/api/auth/sessions", a.GetUserSessionsHandler)

}
