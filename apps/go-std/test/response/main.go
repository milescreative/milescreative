package main

import (
	"fmt"
	"log"
	"net/http"

	"go-std/internal/utils"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `
   <!DOCTYPE html>
   <html>
   <head><title>RESPONSE TEST</title></head>
   <body>
    <h1>CSP Applied</h1>
    <p>Check console. Policy depends on environment.</p>
   </body>
   </html>`)
}

func main() {

	port := "3000"
	mux := http.NewServeMux()
	mux.HandleFunc("/", simpleHandler)
	mux.HandleFunc("/response", responseHandler)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func responseHandler(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, "Hello, World!")
	return

	log.Println("Response sent")

	utils.BadRequest(w, "Bad Request")
}
