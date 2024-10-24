package main

import (
	"gin-pulsar-websockets/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWebSocket(c.Writer, c.Request)
	})

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run Gin server: %v", err)
	}
}
