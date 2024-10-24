package websocket

import (
	"gin-pulsar-websockets/pulsar"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket client connected")

	pulsarClient := pulsar.NewPulsarClient()

	// Listening for messages from WebSocket client
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Send the message to Apache Pulsar
		if err := pulsarClient.Produce(msg); err != nil {
			log.Println("Failed to send message to Pulsar:", err)
		}

		// Optionally echo the message back to WebSocket client
		if err := conn.WriteMessage(messageType, msg); err != nil {
			log.Println("Error writing to WebSocket:", err)
			break
		}
	}
}
