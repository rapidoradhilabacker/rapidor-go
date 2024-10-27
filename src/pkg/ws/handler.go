package ws

// core packages are imported here
// ws is a package that contains the application's websocket logic
// ws is imported in other packages

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager manages WebSocket connections
type ConnectionManager struct {
	connections   map[string]*websocket.Conn
	pulsarManager *PulsarManager
	mu            sync.Mutex
}

// getWSKey returns a unique key for the WebSocket connection
func (cm *ConnectionManager) getWSKey(wsParams WSParams) string {
	if wsParams.TransactionID == "" {
		return "event:" + wsParams.Username + ":" + wsParams.Hostname
	}
	return "event:" + wsParams.Username + ":" + wsParams.Hostname + ":" + wsParams.TransactionID
}

// NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(pm *PulsarManager) *ConnectionManager {
	return &ConnectionManager{
		connections:   make(map[string]*websocket.Conn),
		pulsarManager: pm,
	}
}

// Connect adds a new connection for the user
func (cm *ConnectionManager) Connect(wsConn *websocket.Conn, wsParams WSParams) {
	wsKey := cm.getWSKey(wsParams)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[wsKey] = wsConn
	log.Printf("User connected: %s", wsParams.Username)

	// Start a consumer for this new connection
	go cm.consumeMessages(wsParams, wsConn)
}

// HandleMessages processes messages for the user
func (cm *ConnectionManager) HandleMessages(wsConn *websocket.Conn, wsParams WSParams) {
	wsKey := cm.getWSKey(wsParams)
	defer func() {
		wsConn.Close()
		cm.mu.Lock()
		delete(cm.connections, wsKey)
		cm.mu.Unlock()
		log.Printf("User disconnected: %s", wsParams.Username)
	}()

	for {
		messageType, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Message received from %s: %s", wsParams.Username, string(msg))
		if messageType == websocket.BinaryMessage {
			log.Printf("Received Binary Message from %s: %s", wsParams.Username, string(msg))
			break
		}

		req, err := ParseNotificationFetchRequest(msg)
		if err != nil {
			log.Printf("Error parsing request: %v", err)
			continue
		}

		switch req.ActionType {
		case SYNC_DELTA:
			log.Printf("Received SYNC_DELTA message from %s", wsParams.Username)
		case ACK_RESPONSE:
			log.Printf("Received ACK_RESPONSE message from %s", wsParams.Username)
		case GET_ACTIVE_WS:
			log.Printf("Received GET_ACTIVE_WS message from %s", wsParams.Username)
		case MILESTONE_PROGRESS_UPDATE:
			log.Printf("Received MILESTONE_PROGRESS_UPDATE message from %s", wsParams.Username)
		default:
			log.Printf("Unknown message type: %s", req.ActionType)
		}
	}
}

// consumeMessages listens for messages from Pulsar and sends them to the WebSocket connection
func (cm *ConnectionManager) consumeMessages(wsParams WSParams, conn *websocket.Conn) {
	wsKey := cm.getWSKey(wsParams)
	consumer, err := cm.pulsarManager.Subscribe(wsKey)

	if err != nil {
		log.Printf("Error creating consumer for %s: %v", wsParams.Username, err)
		return
	}
	defer consumer.Close()

	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Printf("Error receiving message for %s: %v", wsParams.Username, err)
			break
		}

		log.Printf("Message received from Pulsar for %s: %s", wsParams.Username, string(msg.Payload()))
		err = conn.WriteMessage(websocket.TextMessage, msg.Payload())
		if err != nil {
			log.Printf("Error sending message to %s: %v", wsParams.Username, err)
			break
		}

		consumer.Ack(msg) // Acknowledge the message after sending
	}
}
