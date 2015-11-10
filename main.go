package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Handler for index.html, shows main view
func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocket handler.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// "Upgrade" HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Register chat session in session registry
	c := &connection{send: make(chan []byte, 256), cs: newChatSession()}
	c.cs.addConnection(c)
	defer c.cs.removeConnection(c)
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writer(&wg, conn)
	go c.reader(&wg, conn)
	wg.Wait()
	conn.Close()
}

func main() {
	initDictionary()
	// Template
	tpl := template.Must(template.ParseFiles("index.html"))

	// Router
	router := http.NewServeMux()
	router.Handle("/", homeHandler(tpl))
	router.HandleFunc("/ws", wsHandler)
	log.Printf("serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
