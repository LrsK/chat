package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type connection struct {
	// Buffered channel of outbound messages.
	send chan []byte
	// chat session
	cs *chatSession
}

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		c.cs.broadcast <- message
	}
}

func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for message := range c.send {
		err := wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

// Handler for index.html, shows main view
func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

type chatSession struct {
	// the mutex to protect connections
	sync.RWMutex
	// Registered connections.
	connections map[*connection]struct{}
	// Inbound messages from the connections.
	broadcast chan []byte
	// Store chat history
	history [][]byte
}

func newChatSession() *chatSession {
	cs := &chatSession{
		broadcast:   make(chan []byte),
		connections: make(map[*connection]struct{}),
		history:     make([][]byte, 0),
	}

	go func() {
		for {
			msg := <-cs.broadcast
			cs.RLock()
			// Send their own message back
			for c := range cs.connections {
				select {
				case c.send <- msg:
				// stop trying to send to this connection after trying for 1 second.
				// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(1 * time.Second):
					log.Printf("shutting down connection %s", c)
					cs.removeConnection(c)
				}
			}
			cs.RUnlock()

			time.Sleep(2 * time.Second)

			// Store message
			cs.history = append(cs.history, msg)
			// Message received, generate response
			response := generateReply(cs.history)
			// Send reply to all other members
			cs.RLock()
			for c := range cs.connections {
				select {
				case c.send <- response:
				// stop trying to send to this connection after trying for 1 second.
				// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(1 * time.Second):
					log.Printf("shutting down connection %s", c)
					cs.removeConnection(c)
				}
			}
			cs.RUnlock()
		}
	}()
	return cs
}

func (cs *chatSession) addConnection(conn *connection) {
	cs.Lock()
	defer cs.Unlock()
	cs.connections[conn] = struct{}{}
	log.Printf("New connection added: %+v", conn)
}

func (cs *chatSession) removeConnection(conn *connection) {
	cs.Lock()
	defer cs.Unlock()
	if _, ok := cs.connections[conn]; ok {
		delete(cs.connections, conn)
		close(conn.send)
	}
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
	// Template
	tpl := template.Must(template.ParseFiles("index.html"))

	// Router
	router := http.NewServeMux()
	router.Handle("/", homeHandler(tpl))
	router.HandleFunc("/ws", wsHandler)
	log.Printf("serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
