package web

import (
	"BarcodeQuery/db"
	model2 "BarcodeQuery/model"
	"github.com/gorilla/websocket"
	"github.com/textileio/go-threads/broadcast"
	"log"
)

type ClientHandlerImpl struct {
	socket     *websocket.Conn
	dbListener *broadcast.Listener
}

func (handler *ClientHandlerImpl) handleQueryDBCallback(queryResult db.DBQueryResult) {
	handler.socket.WriteJSON(model2.BarcodeQueryMessage{
		MessageType: model2.DBQueryNoti,
		Payload:     queryResult,
	})
}

func (handler *ClientHandlerImpl) handle() {

	go func() {
		for {
			v := <-handler.dbListener.Channel()
			msg := v.(db.DBQueryResult)
			handler.handleQueryDBCallback(msg)
		}
	}()1

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
