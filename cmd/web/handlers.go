package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

// mimicks a priliminary db
var snippets = []map[string]interface{} {
	{"id": 21, "text": "learning Go is fun."},
	{"id": 27, "text": "Go Go Go"},
}

// Define a home handler function which writes a byte slice
// containing "Hello from Snippetbox" as the response body.
// NOTE: 
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go") // customize header

	// Initialize a slice containing the paths to the `tmpl` files.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",

	}

	// Read the template file into a template set.
	ts, err := template.ParseFiles(files...)
	// Because the `home()` handler is a method against the application struct
	// it can access its fields, including the structured logger.
	// if err != nil {
	// 	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError) // status code 500
	// 	return
	// }
	if err != nil {
		app.serverError(w, r, err)
	}

	// Write the template content as the response body.
	// As we don't have any `dynamic data` to pass, we assign `nil` to the last parameter.
	// err = ts.Execute(w, nil)
	err = ts.ExecuteTemplate(w, "base", nil) // `base` then invokes `title` & `main` templates.
	// if err != nil {
	// 	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
	if err != nil {
		app.serverError(w, r, err)
	}
	
	// w.Write([]byte("Hello from Snippetbox\n"))
}

// Change the signature of the handler so it's defined as a method against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// Extract the value of the id wildcard & convert it to integer if possible.
	idStr := r.URL.Path[len("/snippet/view/"):] // extract the id from the URL
	id, err := strconv.Atoi(idStr)
	// id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Find the snippet with the matching id:
	var snippet map[string]interface{}
	for _, s := range snippets {
		if s["id"].(int) == id {
			snippet = s
			break
		}
	}

	if snippet == nil {
		http.NotFound(w, r)
		return
	}

	// msg := fmt.Sprintf("Display a specific snippet with ID: %d\n", id)
	// w.Write([]byte(msg))
	// fmt.Fprintf(w, "Display a specific snippet with ID: %d\n", id)
	
	w.Header().Set("Content-Type", "application/json")
	// fmt.Fprintf(w, `{"id", "%d", "name": "Lucy"}`, id)

	json.NewEncoder(w).Encode(snippet)
}


func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet...\n"))
}


func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // 201 status code
	w.Write([]byte("Save a new snippet...\n"))
}
