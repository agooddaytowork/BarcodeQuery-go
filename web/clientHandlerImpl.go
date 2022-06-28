package web

import (
	"github.com/gorilla/websocket"
	"log"
)

type ClientHandlerImpl struct {
}

func (handler *ClientHandlerImpl) handle(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	defer c.Close()
}
