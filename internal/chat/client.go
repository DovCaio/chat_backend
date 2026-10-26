package chat

import (
	"context"
	"encoding/json"
	"log"

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
        send: make(chan []byte, 256),
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
            log.Printf("error when try to get the user connection: %v", err)
            return
        }

        var incoming IncomingMessage
        if err := json.Unmarshal(data, &incoming); err != nil {
            log.Printf("error when try to unmarshal: %v", err)
            continue
        }

        if incoming.To == "" || incoming.Content == "" {
            log.Printf("json empty")
            continue
        }

        outgoing, err := json.Marshal(OutgoingMessage{
            From:    c.userID,
            Content: incoming.Content,
        })

        if err != nil {
            log.Printf("error when try to marshal: %v", err)
            continue
        }
        hub.SendTo(incoming.To, outgoing)
    }
}
