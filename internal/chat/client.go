package chat

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
)

type Client struct {
    userID string
    conn   *websocket.Conn
    send   chan []byte
}

func NewClient(userID string, conn *websocket.Conn ) *Client {
    return &Client{
        userID: userID,
        conn: conn,
    }
}


func (c *Client) WritePump(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return

        case message, ok := <-c.send:
            if !ok {
                return
            }

            if err := c.conn.Write(
                ctx,
                websocket.MessageText,
                message,
            ); err != nil {
                return
            }
        }
    }
}

func (c *Client) ReadPump(
    ctx context.Context,	
	hub *Hub,
) {
    defer hub.Unregister(c)

    for {
        _, data, err := c.conn.Read(ctx)
        if err != nil {
            return
        }

        var incoming IncomingMessage
        if err := json.Unmarshal(data, &incoming); err != nil {
            continue
        }

        if incoming.To == "" || incoming.Content == "" {
            continue
        }

        outgoing, err := json.Marshal(OutgoingMessage{
            From:    c.userID,
            Content: incoming.Content,
        })
        if err != nil {
            continue
        }

        hub.SendTo(incoming.To, outgoing)
    }
}
