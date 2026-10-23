package httpadapter

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"

	"github.com/DovCaio/chat_backend/internal/chat"
)

type WebSocketHandler struct {
    hub *chat.Hub
}

func NewWebSocketHandler(hub *chat.Hub) *WebSocketHandler {
    return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
) {
    userID := r.URL.Query().Get("user")
    if userID == "" {
        http.Error(w, "Usuário obrigatório", http.StatusBadRequest)
        return
    }

    conn, err := websocket.Accept(w, r, nil)
    if err != nil {
        log.Printf("Falha no handshake WebSocket: %v", err)
        return
    }

    ctx, cancel := context.WithCancel(r.Context())
    defer cancel()
    defer conn.CloseNow()

    client := chat.NewClient(userID, conn)

    h.hub.Register(client)

    go client.WritePump(ctx)
    client.ReadPump(ctx, h.hub)
}
