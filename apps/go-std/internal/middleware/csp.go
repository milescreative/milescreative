package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// CSPMiddleware: Content Security Policy and security headers, adapts to env.
func CSPMiddleware(isDevelopment bool) func(http.HandlerFunc) http.HandlerFunc {
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

	return func(next http.HandlerFunc) http.HandlerFunc {
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

			next(w, r)
		})
	}
}
