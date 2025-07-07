package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	serveMux.HandleFunc("/healthz", handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(srv.ListenAndServe())

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
