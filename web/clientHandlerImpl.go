package web

import (
	model2 "BarcodeQuery/model"
	"github.com/gorilla/websocket"
	"github.com/textileio/go-threads/broadcast"
	"log"
)

type ClientHandlerImpl struct {
	socket     *websocket.Conn
	dbListener *broadcast.Listener
}

func (handler *ClientHandlerImpl) handleMessageCB(msg model2.BarcodeQueryMessage) {
	handler.socket.WriteJSON(msg)
}

func (handler *ClientHandlerImpl) handle() {
	go func() {
		for {
			v := <-handler.dbListener.Channel()
			msg := v.(model2.BarcodeQueryMessage)
			handler.handleMessageCB(msg)
		}
	}()
	for {
		mt, message, err := handler.socket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println(mt, message)
	}
	defer handler.socket.Close()
}
