package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Pool struct {
    Register   chan *Client
    Unregister chan *Client
    Clients    map[*Client]bool
    Broadcast  chan BroadcastMessage  // Ensure this uses Message struct
}

type BroadcastMessage struct {
    Sender  *Client
    Message Message
}

func NewPool() *Pool {
    return &Pool{
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan BroadcastMessage),
    }
}

// Function to handle incoming messages and broadcast them
func (pool *Pool) Start() {
    for {
        select {
				// TODO: Implement the Register case and the user name
        case client := <-pool.Register:
            pool.Clients[client] = true
            log.Printf("Size of Connection Pool: %d", len(pool.Clients))
            // Broadcasting a system message that a client has joined
            // pool.broadcastMessage(
						// 		nil,
						// 		Message{User: "System", Body: "Someone has joined the chat"},
						// )

        case client := <-pool.Unregister:
            if _, exists := pool.Clients[client]; exists {
                delete(pool.Clients, client)
                client.Conn.Close()
                // Broadcasting a system message that a client has left
                pool.broadcastMessage(
									nil,
									Message{User: "System", Body: "Someone has left the chat"},
							)
            }

        case broadcast := <-pool.Broadcast:
            fmt.Println("Message Received: ", broadcast.Message)
            pool.broadcastMessage(broadcast.Sender, broadcast.Message)
        }
    }
}

// Broadcasts a message to all clients except the sender
func (pool *Pool) broadcastMessage(sender *Client, message Message) {
    fmt.Println("Broadcasting Message: ", message.Body)
    jsonMsg, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshalling message: %v", err)
        return
    }

    for client := range pool.Clients {
        if sender == nil || client.ID != sender.ID {
            err := client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
            if err != nil {
                delete(pool.Clients, client)
                client.Conn.Close()
                log.Printf("Error sending message to client %s: %v", client.ID, err)
            }
        }
    }
}
