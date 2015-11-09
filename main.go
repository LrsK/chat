package main

import (
	"html/template"
	"log"
	"net/http"
)

func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

// Set up routers
// Set up handler
// Listen and serve

func main() {

	// Template
	tpl := template.Must(template.ParseFiles("index.html"))

	// Router
	router := http.NewServeMux()
	router.Handle("/", homeHandler(tpl))
	router.Handle("/ws", wsHandler)
	log.Printf("serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
