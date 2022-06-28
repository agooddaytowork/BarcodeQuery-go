package web

import "github.com/gorilla/websocket"

type ClientHandler interface {
	handle(c *websocket.Conn)
}
