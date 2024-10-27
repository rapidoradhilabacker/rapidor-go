package ws

// core packages are imported here
// ws is a package that contains the application's websocket logic
// ws is imported in other packages

import (
	"go-pulsar-websockets/src/core"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Create an instance of websocket.Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SetupWebSocketRoutes sets up the WebSocket routes
func SetupWebSocketRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1/ws")
	{
		v1.GET("/", handleWebSocketConnection)
	}
}

// handleWebSocketConnection upgrades the HTTP connection to a WebSocket connection
func handleWebSocketConnection(c *gin.Context) {
	env := core.LoadEnv()

	pulsarManager, err := NewPulsarManager(env.PulsarURL)
	if err != nil {
		log.Fatalf("Failed to create Pulsar manager: %v", err)
	}
	defer pulsarManager.Close()

	// Create a new connection manager
	cm := NewConnectionManager(pulsarManager)

	// Upgrade HTTP connection to WebSocket
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	wsParams := WSParams{
		Username:      c.Query("username"),
		Hostname:      c.Query("hostname"),
		TransactionID: c.Query("transaction_id"),
	}

	// Connect the WebSocket with the connection manager
	cm.Connect(wsConn, wsParams)
	go cm.HandleMessages(wsConn, wsParams)
}
