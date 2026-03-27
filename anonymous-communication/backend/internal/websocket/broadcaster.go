package websocket

func BroadcastToUser(hub *Hub, userID string, payload any) bool {
	if hub == nil {
		return false
	}

	return hub.SendToUser(userID, payload)
}
