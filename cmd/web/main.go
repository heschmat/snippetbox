package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	// Import the models package.
	"github.com/heschmat/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Hold the applicatin-wide dependencies for the web app.
type application struct {
	logger *slog.Logger
	// This makes the SnippetModel object available to the handlers:
	snippets *models.SnippetModel
}

func main() {
	// Define a new command-line flag with the name `addr`
	// flag.String() return a pointer to the flag value.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for MySQL DSN string.
	dsn := flag.String("dsn", "web:summer2006@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// WE pass openDB() the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer call to `db.Close()` so that the connection pool is closed before `main()` function exits.
	defer db.Close()

	// Initialize a new instance of the application struct, containing the dependencies.
	app := &application{
		logger: logger,
		// Initialize a models.SnippetModel instance containing the connection pool.
		snippets: &models.SnippetModel{DB: db},
	}

	mux := http.NewServeMux() // Initialize a new servemux

	fileServer := http.FileServer(http.Dir("./ui/static"))
	// the following gives you `404 page not found`
	// when visting: http://localhost:4000/static/
	// mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))


	// Register the `home` func as the handler for the "/" URL pattern.
	// Restrict this route to exact matches on "/" only, using {$}
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate) // shows the form
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost) //

	// log.Printf("Starting server on %s", *addr)
	logger.Info("starting server", "addr", *addr)

	// now we can specify the port manually: 
	// $ go run ./cmd/web -addr=":3000"
	err = http.ListenAndServe(*addr, mux)
	// log.Fatal(err) // Log the error message & exit. 
	// N.B. Any error returned by `http.ListenAndServe()` is always `non-nil`

	logger.Error(err.Error())
	os.Exit(1)
}

// This function wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	// Initialize the connection pool for future use.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// N.B. database connections are established.
	// Here, we ping the connection to verify if everything is set up correct
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
