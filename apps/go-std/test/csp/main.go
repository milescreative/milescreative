package main

import (
	"fmt"
	"go-std/internal/config"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// CSPMiddleware: Content Security Policy and security headers, adapts to env.
func CSPMiddleware(isDevelopment bool) func(http.Handler) http.Handler {
	prodCspDirectives := map[string][]string{
		"default-src": {"'none'"},
		"script-src":  {"'self'"},
		"style-src": {
			"'self'",
			"fonts.googleapis.com",
		},
		"img-src": {"'self'", "data:"},
		"font-src": {
			"'self'",
			"fonts.gstatic.com",
			"fonts.googleapis.com",
		},
		"connect-src":               {"'self'"},
		"media-src":                 {"'self'"},
		"object-src":                {"'none'"},
		"base-uri":                  {"'self'"},
		"form-action":               {"'self'"},
		"frame-ancestors":           {"'none'"},
		"upgrade-insecure-requests": {},
		"block-all-mixed-content":   {},
		"frame-src":                 {"'none'"},
		"manifest-src":              {"'self'"},
		"worker-src":                {"'self'"},
	}

	devCspDirectives := map[string][]string{
		"default-src": {"'self'"},
		"script-src":  {"'self'", "'unsafe-inline'", "'unsafe-eval'"},
		"style-src":   {"'self'", "'unsafe-inline'"},
		"img-src":     {"'self'", "data:"},
		"font-src":    {"'self'"},
		"connect-src": {"'self'", "ws://localhost:3000"}, // Limit to dev server
		"form-action": {"'self'"},
		"frame-src":   {"'none'"},
		"object-src":  {"'none'"},
		"base-uri":    {"'self'"},
	}

	var selectedDirectives map[string][]string
	if isDevelopment {
		selectedDirectives = devCspDirectives
		log.Println("CSP: DEVELOPMENT policy")
	} else {
		selectedDirectives = prodCspDirectives
		log.Println("CSP: PRODUCTION policy")
	}

	var policyParts []string
	for directive, sources := range selectedDirectives {
		if len(sources) > 0 {
			policyParts = append(policyParts, fmt.Sprintf("%s %s", directive, strings.Join(sources, " ")))
		} else {
			policyParts = append(policyParts, directive)
		}
	}

	if !isDevelopment {
		policyParts = append(policyParts, "report-uri /csp-report") // Important!
	}
	cspValue := strings.Join(policyParts, "; ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", cspValue)
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("X-XSS-Protection", "0")

			if !isDevelopment && (r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https") {
				maxAgeSeconds := 31536000
				hstsValue := fmt.Sprintf("max-age=%d; includeSubDomains", maxAgeSeconds)
				w.Header().Set("Strict-Transport-Security", hstsValue)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `
   <!DOCTYPE html>
   <html>
   <head><title>Conditional CSP Test</title></head>
   <body>
    <h1>CSP Applied</h1>
    <p>Check console. Policy depends on environment.</p>
   </body>
   </html>`)
}

func main() {
	env, _ := config.Config()
	isDev := strings.ToLower(env.GetString("APP_ENV")) == "development"

	cspMiddlewareInstance := CSPMiddleware(isDev)
	mux := http.NewServeMux()
	mux.HandleFunc("/", simpleHandler)
	secureMux := cspMiddlewareInstance(mux)

	try_env := loadConfig()
	fmt.Println("app_name: ", try_env.GetString("app_name"))

	port := strconv.Itoa(try_env.Port())
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting server on port %s (Dev Mode: %t)...", port, isDev)
	err := http.ListenAndServe(":"+port, secureMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func loadConfig() *config.ConfigMap {
	env, _ := config.Config()
	env.LoadJSON(filepath.Join("config", "test.jsonc"))
	return env
}
