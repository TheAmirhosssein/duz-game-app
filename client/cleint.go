package client

import "github.com/gorilla/websocket"

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
