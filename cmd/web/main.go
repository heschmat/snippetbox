package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Define a new command-line flag with the name `addr`
	// flag.String() return a pointer to the flag value.
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

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

	// log.Printf("Starting server on %s", *addr)
	logger.Info("starting server", "addr", *addr)

	// now we can specify the port manually: 
	// $ go run ./cmd/web -addr=":3000"
	err := http.ListenAndServe(*addr, mux)
	// log.Fatal(err) // Log the error message & exit. 
	// N.B. Any error returned by `http.ListenAndServe()` is always `non-nil`

	logger.Error(err.Error())
	os.Exit(1)
}
