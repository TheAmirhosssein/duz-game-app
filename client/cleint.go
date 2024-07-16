package client

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	UserId string
}

func New(id, userId string, conn *websocket.Conn) *Client {
	return &Client{
		ID:     id,
		Conn:   conn,
		UserId: userId,
	}
}

func (client Client) SendMessageToClient(message []byte) {
	err := client.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Write error:", err)
		client.Conn.Close()
	}
}
