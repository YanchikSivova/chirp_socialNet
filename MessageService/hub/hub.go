package hub

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type Hub struct {
	clients map[uuid.UUID]map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[uuid.UUID]map[*websocket.Conn]bool)}
}

func (h *Hub) AddClient(profileID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[profileID]; !ok {
		h.clients[profileID] = make(map[*websocket.Conn]bool)
	}
	h.clients[profileID][conn] = true
}
func (h *Hub) RemoveClient(profileID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	connections, ok := h.clients[profileID]
	if !ok {
		return
	}
	delete(connections, conn)
	if len(connections) == 0 {
		delete(h.clients, profileID)
	}
}

func (h *Hub) Send(profileID uuid.UUID, payload any) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	connections, ok := h.clients[profileID]
	if !ok {
		return nil
	}
	for conn := range connections {
		err := conn.WriteJSON(payload)
		if err != nil {
			return conn.Close()
		}
	}
	return nil
}

func (h *Hub) IsOnline(profileID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	connections, ok := h.clients[profileID]
	return ok && len(connections) > 0
}
