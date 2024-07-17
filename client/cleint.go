package client

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID          string
	Conn        *websocket.Conn
	UserId      string
	PawnOnBoard int
}

func New(id, userId string, conn *websocket.Conn) *Client {
	return &Client{
		ID:          id,
		Conn:        conn,
		UserId:      userId,
		PawnOnBoard: 0,
	}
}

func (client Client) SendMessageToClient(message []byte) {
	err := client.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Write error:", err)
		client.Conn.Close()
	}
}

func (client Client) MaxMove() bool {
	return client.PawnOnBoard == 3
}

func (client *Client) MovedPawn() {
	client.PawnOnBoard++
}
