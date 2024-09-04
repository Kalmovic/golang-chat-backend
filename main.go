package main

import (
	"fmt"
	"net/http"

	"github.com/Kalmovic/golang-chat/pkg/websocket"
	"github.com/google/uuid"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	client := &websocket.Client{
		// Generate a unique ID for the client
		ID: uuid.NewString(),
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	/**
		* go keyword is used to create a new fiber that will run concurrently with the rest of the program.
	*/
	go pool.Start();

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes();
	http.ListenAndServe(":8080", nil)
}