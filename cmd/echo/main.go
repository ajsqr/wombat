package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Hello from my mac")
	})

	addr := "127.0.0.1:1233"
	log.Printf("HTTP server listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
