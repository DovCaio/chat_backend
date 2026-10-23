package chat

import (
	"sync"
)


type Hub struct {
    mu      sync.Mutex
    clients map[string]map[*Client]bool
}

func NewHub() *Hub {
    return &Hub{
        clients: make(map[string]map[*Client]bool),
    }
}

func (h *Hub) Register(c *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if h.clients[c.userID] == nil {
        h.clients[c.userID] = make(map[*Client]bool)
    }

    h.clients[c.userID][c] = true
}

func (h *Hub) Unregister(c *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if userClients := h.clients[c.userID]; userClients != nil {
        if userClients[c] {
            delete(userClients, c)
            close(c.send)
        }

        if len(userClients) == 0 {
            delete(h.clients, c.userID)
        }
    }
}


func (h *Hub) SendTo(userID string, message []byte) {
    h.mu.Lock()
    defer h.mu.Unlock()

    for client := range h.clients[userID] {
        select {
        case client.send <- message:
        default:
        }
    }
}
