package main

import (
	"fmt"
	"go-std/internal/auth"
	"go-std/internal/config"
	"go-std/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	env, _ := config.Config()

	isDev := strings.ToLower(env.GetString("APP_ENV")) == "development"

	mux := http.NewServeMux()

	mux.HandleFunc("/", envHandler)
	mux.HandleFunc("/login", simpleHandler)
	mux.HandleFunc("/api/auth/callback/google", callbackHandler)
	port, _ := env.GetInt("APP_PORT")
	if port == 0 {
		port = 3000
	}
	portStr := strconv.Itoa(port)

	log.Printf("Starting server on port %s (Dev Mode: %t)...", portStr, isDev)
	err := http.ListenAndServe(":"+portStr, mux)
	if err != nil {
		log.Fatal(err)
	}

}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateState()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("state:", state)

	codeVerifier, err := utils.GenerateCodeVerifier()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("codeVerifier:", codeVerifier)
	env, _ := config.Config()

	isDev := strings.ToLower(env.GetString("APP_ENV")) == "development"

	goog := auth.NewGoogleOAuth(
		env.GetString("GOOGLE_CLIENT_ID"),
		env.GetString("GOOGLE_CLIENT_SECRET"),
		env.GetString("GOOGLE_REDIRECT_URI"),
		[]string{"email", "profile"},
	)

	u, err := goog.CreateAuthorizationURLWithPKCE(state, codeVerifier)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("u:", u)

	cookies := []*http.Cookie{{
		Name:     "google_oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteLaxMode,
	},
		{
			Name:     "google_code_verifier",
			Value:    codeVerifier,
			Path:     "/",
			HttpOnly: true,
			Secure:   !isDev,
			SameSite: http.SameSiteLaxMode,
		},
	}
	for _, c := range cookies {
		http.SetCookie(w, c)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `
   <!DOCTYPE html>
   <html>
   <head><title>Conditional CSP Test</title></head>
   <body>
    <h1>OAuth 2.0 Authentication Example</h1>
    <p>Click the button below to sign in with Google:</p>
    <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #4285f4; color: white; text-decoration: none; border-radius: 4px;">Sign in with Google</a>
   </body>
   </html>`, u)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	env, _ := config.Config()
	//get state from cookie
	storedState, err := r.Cookie("google_oauth_state")
	if err != nil {
		log.Fatal(err)
	}
	storedCodeVerifier, err := r.Cookie("google_code_verifier")
	if err != nil {
		log.Fatal(err)
	}

	if storedState.Value != state {
		log.Fatal("state mismatch- please restart")
	}
	goog := auth.NewGoogleOAuth(
		env.GetString("GOOGLE_CLIENT_ID"),
		env.GetString("GOOGLE_CLIENT_SECRET"),
		env.GetString("GOOGLE_REDIRECT_URI"),
		[]string{"email", "profile"},
	)
	tokens, err := goog.ValidateAuthorizationCode(code, storedCodeVerifier.Value)
	if err != nil {
		log.Fatal(err)
	}
	idToken, err := tokens.IDToken()
	if err != nil {
		log.Fatal(err)
	}

	claims, err := utils.DecodeJwt(idToken)
	if err != nil {
		log.Fatal(err)
	}

	googleId := claims["sub"].(string)
	name := claims["name"].(string)
	picture := claims["picture"].(string)
	email := claims["email"].(string)
	verified := claims["email_verified"].(bool)

	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
        <head><title>Callback Page</title></head>
        <body>
            <h1>Callback Page</h1>
            <p>This is the callback page.</p>
						<p>Google ID: %s</p>
						<p>Name: %s</p>
						<p>Picture: %s</p>
						<p>Email: %s</p>
						<p>Verified: %t</p>
        </body>
        </html>
    `, googleId, name, picture, email, verified)
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	env, _ := config.Config()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
        <head><title>Environment Variables</title></head>
        <body>
            <h1>Environment Variables</h1>
            <pre>
APP_ENV: %s
GOOGLE_CLIENT_ID: %s
GOOGLE_CLIENT_SECRET: %s
GOOGLE_REDIRECT_URI: %s
PORT: %d
            </pre>
        </body>
        </html>
    `, env.GetString("APP_ENV"),
		env.GetString("GOOGLE_CLIENT_ID"),
		env.GetString("GOOGLE_CLIENT_SECRET"),
		env.GetString("GOOGLE_REDIRECT_URI"),
		env.Port())
}
