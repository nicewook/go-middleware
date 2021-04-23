package main

import (
	"fmt"
	"log"
	"net/http"
)

func mOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start mOne")
		next.ServeHTTP(w, r)
		log.Println("--end mOne")
	})
}

func mTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start mTwo")
		if r.URL.Path == "/foo" {
			msg := fmt.Sprint(r.URL.Path, ": so not running next middleWare")
			log.Println(msg)
			w.Write([]byte(msg))
			return
		}
		next.ServeHTTP(w, r)
		log.Println("--end mTwo")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("finalHandler")
	w.Write([]byte("finalHandler"))
}

func main() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)

	mux.Handle("/", mOne(mTwo(finalHandler)))

	port := "3000"
	log.Println("listening on: ", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
