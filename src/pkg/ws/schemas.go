package ws

// core packages are imported here
// ws is a package that contains the application's websocket logic
// ws is imported in other packages

import (
	"encoding/json"
	"fmt"
)

type WSParams struct {
	Username      string
	Hostname      string
	TransactionID string
}

type NotificationActionType string

const (
	SYNC_DELTA                NotificationActionType = "sync_delta"
	ACK_RESPONSE              NotificationActionType = "ack_response"
	GET_ACTIVE_WS             NotificationActionType = "get_active_ws"
	MILESTONE_PROGRESS_UPDATE NotificationActionType = "milestone_progress_update"
)

func (action NotificationActionType) String() string {
	return string(action)
}

type NotificationFetchRequest struct {
	ActionType NotificationActionType `json:"action_type"`
	Hostname   string                 `json:"hostname"`
	Username   string                 `json:"username"`
	UUID       string                 `json:"uuid"`
}

// ParseNotificationFetchRequest parses JSON message into NotificationFetchRequest struct
func ParseNotificationFetchRequest(message []byte) (*NotificationFetchRequest, error) {
	var req NotificationFetchRequest
	err := json.Unmarshal(message, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification request: %v", err)
	}
	return &req, nil
}
