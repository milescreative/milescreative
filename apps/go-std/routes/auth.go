package routes

import (
	"go-std/internal/auth"
	"go-std/internal/config"
	"go-std/internal/utils"
	"log"
	"net/http"

	"github.com/g-h-miles/httpmux"
)

func AuthRoutes(mux *httpmux.Router, app *config.App) {

	r := mux

	a, err := auth.NewAuthHandlers(app)
	if err != nil {
		log.Fatalf("failed to create auth handlers: %v", err)
	}

	//todo: move to root
	r.GET("/{$}", DummyHandler)
	r.GET("/api/auth/login", a.LoginHandler)
	r.GET("/api/auth/login/{provider}", a.LoginHandler)
	r.GET("/api/auth/callback", a.CallbackHandler)
	r.GET("/api/auth/callback/{provider}", a.CallbackHandler)
	r.GET("/api/auth/validate", a.ValidateSessionHandler)
	r.POST("/api/auth/refresh", a.RefreshTokenHandler)
	r.GET("/api/auth/logout", a.LogoutHandler) //todo: change to POST
	r.GET("/api/auth/csrf", a.GetCSRFTokenHandler)
	r.GET("/api/auth/user", a.GetUserHandler)
	r.POST("/api/auth/user", a.UpdateUserHandler)
	r.DELETE("/api/auth/user", a.DeleteUserHandler)
	r.GET("/api/auth/sessions", a.GetUserSessionsHandler)

}

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, "dummy handler")
}

func AuthRoutesStd(mux *httpmux.Router, app *config.App) {

	r := mux

	a, err := auth.NewAuthHandlers(app)
	if err != nil {
		log.Fatalf("failed to create auth handlers: %v", err)
	}

	//todo: move to root
	r.GET("/{$}", DummyHandler)
	r.GET("/api/auth/login", a.LoginHandler)
	r.GET("/api/auth/login/{provider}", a.LoginHandler)
	r.GET("/api/auth/callback", a.CallbackHandler)
	r.GET("/api/auth/callback/{provider}", a.CallbackHandler)
	r.GET("/api/auth/validate", a.ValidateSessionHandler)
	r.POST("/api/auth/refresh", a.RefreshTokenHandler)
	r.GET("/api/auth/logout", a.LogoutHandler)
	r.GET("/api/auth/csrf", a.GetCSRFTokenHandler)
	r.GET("/api/auth/user", a.GetUserHandler)
	r.POST("/api/auth/user", a.UpdateUserHandler)
	r.DELETE("/api/auth/user", a.DeleteUserHandler)
	r.GET("/api/auth/sessions", a.GetUserSessionsHandler)

}
