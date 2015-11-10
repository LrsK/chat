package main

import (
	"log"
	"sync"
	"time"
)

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
			// Simulate thinking
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
