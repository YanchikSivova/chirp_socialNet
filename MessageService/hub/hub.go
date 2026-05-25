package hub

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type Hub struct {
	clients map[uuid.UUID]*websocket.Conn
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[uuid.UUID]*websocket.Conn)}
}

func (h *Hub) AddClient(profileID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[profileID] = conn
}
func (h *Hub) RemoveClient(profileID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, profileID)
}

func (h *Hub) Send(profileID uuid.UUID, payload any) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.clients[profileID]
	if !ok {
		return nil
	}
	return conn.WriteJSON(payload)
}
