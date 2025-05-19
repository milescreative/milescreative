package routes

import (
	"go-std/internal/auth"
	"go-std/internal/config"
	"net/http"
)

func AuthRoutes(mux *http.ServeMux, app *config.App) {

	router := NewRouter(mux)
	r := router.ROUTE

	a := auth.NewAuthHandlers(app)

	//todo: move to root
	r("/").GET(ROOT(a.EnvHandler))

	r("/api/auth/login/google").GET(a.LoginHandler)
	r("/api/auth/callback/google").GET(a.CallbackHandler)
	r("/api/auth/validate").GET(a.ValidateSessionHandler)
	r("/api/auth/refresh").POST(a.RefreshTokenHandler)
	r("/api/auth/logout").POST(a.LogoutHandler)
	r("/api/auth/csrf").GET(a.GetCSRFTokenHandler)
	r("/api/auth/user").
		GET(a.GetUserHandler).
		POST(a.UpdateUserHandler).
		DELETE(a.DeleteUserHandler)
	r("/api/auth/sessions").GET(a.GetUserSessionsHandler)

}
