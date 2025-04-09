package main

import (
	"encoding/json"
	"go-std/internal/middleware"
	ratelimit "go-std/internal/middleware/rate-limit"
	"go-std/internal/utils"
	"log"
	"net/http"
)

func main() {
	sessionId, err := utils.GenerateRandomString()
	if err != nil {
		log.Fatal(err)
	}
	csrfToken, csrfTokenHMAC := utils.GenerateCSRFToken(sessionId)
	log.Println("csrfToken: ", csrfToken)
	log.Println("csrfTokenHMAC: ", csrfTokenHMAC)
	code, err := utils.GenerateRandomString()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(code)

	valid := utils.VerifyCSRFToken(csrfToken, sessionId, csrfTokenHMAC)
	log.Println("valid: ", valid)
	router := http.NewServeMux()
	router.HandleFunc("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Write([]byte("received request for item: " + id))
	})

	router.HandleFunc("POST /postTest", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("received post request"))
	})

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.TokenBucket(5, 1, "test"),
		middleware.CORS,
	)

	router.HandleFunc("GET /test-cors", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("CORS test successful"))
	})

	router.HandleFunc("GET /test-rate", func(w http.ResponseWriter, r *http.Request) {
		storeSize := ratelimit.GetStoreSize()
		log.Printf("Rate limit test hit: %s, Store size: %d", r.URL.Path, storeSize)

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message":    "Rate limit test successful",
			"store_size": storeSize,
			"store_keys": ratelimit.GetStoreKeys(),
		}
		json.NewEncoder(w).Encode(response)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	log.Println("Starting server on port 8080")
	server.ListenAndServe()

}
