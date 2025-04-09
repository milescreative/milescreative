package middleware

import (
	ratelimit "go-std/internal/middleware/rate-limit"
	"log"
	"net/http"
)

// TokenBucket returns a middleware that uses the token bucket algorithm for rate limiting
func TokenBucket(bucketCapacity int, refillRate float64, prefix string) func(http.Handler) http.Handler {
	limiter := ratelimit.NewTokenBucketRateLimiter(bucketCapacity, refillRate, prefix)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.IsAllowed(r.URL.Path) {
				log.Printf("Rate limit exceeded for path: %s", r.URL.Path)
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			log.Printf("Rate limit allowed for path: %s", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
