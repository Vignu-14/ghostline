package websocket

import "sync"

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if existing, ok := h.clients[client.UserID]; ok && existing != client {
		_ = existing.Close()
	}

	h.clients[client.UserID] = client
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	existing, ok := h.clients[client.UserID]
	if ok && existing == client {
		delete(h.clients, client.UserID)
	}
}

func (h *Hub) SendToUser(userID string, payload any) bool {
	h.mu.RLock()
	client := h.clients[userID]
	h.mu.RUnlock()

	if client == nil {
		return false
	}

	if err := client.WriteJSON(payload); err != nil {
		h.Unregister(client)
		_ = client.Close()
		return false
	}

	return true
}
