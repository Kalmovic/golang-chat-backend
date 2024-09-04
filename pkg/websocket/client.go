package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single chat client.
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Body struct {
	User string `json:"user"`
	Body string `json:"body"`
}

// Message defines the structure of a message sent over the websocket.
type Message struct {
	User string `json:"user"`
  Body string `json:"body"`
}

// Read continuously reads from the websocket connection and handles messages.
func (c *Client) Read() {
	defer func() {
		// Unregister the client and close the connection when this function exits.
		c.Pool.Unregister <- c
		if err := c.Conn.Close(); err != nil {
			log.Printf("Error closing connection for client %s: %v", c.ID, err)
		}
	}()

	for {
		// Read in a new message.
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client %s: %v", c.ID, err)
			break
		}

		// Convert the payload into a Message struct
		var message Message
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Convert Message struct to JSON for broadcasting
		// jsonMsg, err := json.Marshal(message)
		// if err != nil {
		// 	log.Printf("Error marshalling message: %v", err)
		// 	continue
		// }

		c.Pool.Broadcast <- BroadcastMessage{
			Sender:  c,
			Message: message,
		}

		// Log the received and broadcasted message for debugging
		log.Printf("Message Received and Broadcasted from %s: %+v\n", c.ID, message)
	}
}
