package websocket

import (
	"sync"

	fiberws "github.com/gofiber/contrib/websocket"
)

type Client struct {
	UserID string
	Conn   *fiberws.Conn
	mu     sync.Mutex
}

func NewClient(userID string, conn *fiberws.Conn) *Client {
	return &Client{
		UserID: userID,
		Conn:   conn,
	}
}

func (c *Client) WriteJSON(payload any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.WriteJSON(payload)
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.Close()
}
