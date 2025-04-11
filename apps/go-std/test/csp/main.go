package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // To read environment variables
	"strings"
)

// CSPMiddleware provides Content Security Policy and other security headers.
// Adjusts policy based on development vs. production environment.
// Assumes assets are served from the same origin ('self'), potentially cached by Cloudflare.
func CSPMiddleware(isDevelopment bool) func(http.Handler) http.Handler {
	// --- Define Production CSP Directives ---
	// Start with a strict policy for production
	prodCspDirectives := map[string][]string{
		"default-src":               {"'self'"},
		"script-src":                {"'self'"},          // Production builds should avoid inline/eval
		"style-src":                 {"'self'"},          // Aim for external CSS files. Add 'unsafe-inline' ONLY if absolutely necessary.
		"img-src":                   {"'self'", "data:"}, // Allow self and data URIs
		"font-src":                  {"'self'"},          // Allow self. Add external font providers if needed.
		"connect-src":               {"'self'"},          // API calls to self. Add external APIs if needed.
		"form-action":               {"'self'"},
		"frame-src":                 {"'none'"}, // Disallow framing
		"object-src":                {"'none'"}, // Disallow plugins
		"base-uri":                  {"'self'"},
		"frame-ancestors":           {"'none'"}, // Disallow embedding in frames
		"upgrade-insecure-requests": {},
	}

	// --- Define Development CSP Directives ---
	// Looser policy for development to allow HMR, inline styles/scripts, dev tools etc.
	devCspDirectives := map[string][]string{
		"default-src": {"'self'"},
		// Allow inline scripts and eval for dev servers, HMR etc.
		"script-src": {"'self'", "'unsafe-inline'", "'unsafe-eval'"},
		// Allow inline styles often used during development or by UI libraries.
		"style-src": {"'self'", "'unsafe-inline'"},
		"img-src":   {"'self'", "data:"},
		"font-src":  {"'self'"},
		// Allow connections to self and potentially WebSocket connections for HMR.
		// Use '*' carefully, prefer specific dev server ports if known (e.g., ws://localhost:*)
		"connect-src": {"'self'", "ws:"}, // Allow WebSocket connections
		"form-action": {"'self'"},
		"frame-src":   {"'none'"}, // Keep restrictive unless needed for dev tools
		"object-src":  {"'none'"},
		"base-uri":    {"'self'"},
		// Allow framing by self for certain dev tools if necessary
		// "frame-ancestors": {"'self'"}, // Uncomment if needed
	}
	// Note: upgrade-insecure-requests is omitted in dev as HTTPS might not be used locally.

	// --- Select Policy Based on Environment ---
	var selectedDirectives map[string][]string
	if isDevelopment {
		selectedDirectives = devCspDirectives
		log.Println("CSP Middleware: Applying DEVELOPMENT policy")
	} else {
		selectedDirectives = prodCspDirectives
		log.Println("CSP Middleware: Applying PRODUCTION policy")
	}

	// --- Assemble CSP String ---
	var policyParts []string
	for directive, sources := range selectedDirectives {
		if len(sources) > 0 {
			policyParts = append(policyParts, fmt.Sprintf("%s %s", directive, strings.Join(sources, " ")))
		} else {
			policyParts = append(policyParts, directive) // Directives with no value
		}
	}
	cspValue := strings.Join(policyParts, "; ")

	// --- Return the actual middleware handler ---
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// --- Set Headers ---

			// Set the CSP header. Use Report-Only during initial setup/testing.
			// w.Header().Set("Content-Security-Policy-Report-Only", cspValue)
			w.Header().Set("Content-Security-Policy", cspValue)

			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("X-XSS-Protection", "0")

			// Only set HSTS in production AND if HTTPS is enforced.
			if !isDevelopment {
				// Ensure HTTPS is actually being used before sending HSTS
				// A common check is for the X-Forwarded-Proto header if behind a proxy like Traefik/Cloudflare
				isHTTPS := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
				if isHTTPS {
					maxAgeSeconds := 31536000 // 1 year
					hstsValue := fmt.Sprintf("max-age=%d", maxAgeSeconds)
					// Add includeSubDomains / preload if applicable
					w.Header().Set("Strict-Transport-Security", hstsValue)
				} else {
					// Optional: Log a warning if in production but not serving over HTTPS
					// log.Println("Warning: HSTS header skipped - connection is not secure.")
				}
			}

			// --- Call next handler ---
			next.ServeHTTP(w, r)
		})
	}
}

// --- Example Usage ---

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head><title>Conditional CSP Test</title></head>
		<body>
			<h1>CSP Applied</h1>
			<p>Check developer console for CSP reports/errors. Policy depends on environment.</p>
			<!-- Dev mode might allow inline scripts/styles here, Prod mode would block -->
			<!-- <script>console.log('Inline script test');</script> -->
			<!-- <p style="color: red;">Inline style test</p> -->
		</body>
		</html>`)
}

func main() {
	// --- Determine Environment ---
	// Example: Use an environment variable. Default to production if not set.
	env := os.Getenv("APP_ENV") // e.g., set APP_ENV=development
	isDev := strings.ToLower(env) == "development"

	// --- Create Middleware Instance ---
	cspMiddlewareInstance := CSPMiddleware(isDev)

	// --- Setup Routing ---
	mux := http.NewServeMux()
	mux.HandleFunc("/", simpleHandler)
	// Add your API handlers etc.

	// --- Apply Middleware ---
	secureMux := cspMiddlewareInstance(mux) // Apply the configured CSP middleware

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Printf("Starting server on port %s (Development Mode: %t)...", port, isDev)
	// Use ListenAndServeTLS in production!
	err := http.ListenAndServe(":"+port, secureMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
