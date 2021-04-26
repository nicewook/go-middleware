package main

import (
	"log"
	"net/http"
)

func CheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start CheckAuthentication")
		if r.Header.Get("Client-ID") != "hsjeong" {
			http.Error(w, "Client not found", http.StatusForbidden)
			return
		}
		if r.Header.Get("Client-Access-ID") != "hsjeong-access" {
			http.Error(w, "Invalid secret key", http.StatusForbidden)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Not proper content type", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
		log.Println("--end CheckAuthentication")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("finalHandler")
	w.Write([]byte("finalHandler"))
}

func main() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)

	mux.Handle("/", CheckAuthentication(finalHandler))

	port := "3000"
	log.Println("listening on: ", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
