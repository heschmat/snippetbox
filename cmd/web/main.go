package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux() // Initialize a new servemux

	fileServer := http.FileServer(http.Dir("./ui/static"))
	// the following gives you `404 page not found`
	// when visting: http://localhost:4000/static/
	// mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))


	// Register the `home` func as the handler for the "/" URL pattern.
	// Restrict this route to exact matches on "/" only, using {$}
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate) // shows the form
	mux.HandleFunc("POST /snippet/create", snippetCreatePost) //

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err) // Log the error message & exit. 
	// N.B. Any error returned by `http.ListenAndServe()` is always `non-nil`
}
