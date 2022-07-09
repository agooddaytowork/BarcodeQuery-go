package web

import (
	"BarcodeQuery/db"
	model2 "BarcodeQuery/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/textileio/go-threads/broadcast"
	"log"
)

type ClientHandlerImpl struct {
	socket          *websocket.Conn
	dbListener      *broadcast.Listener
	clientBroadcast *broadcast.Broadcaster
}

func (handler *ClientHandlerImpl) handleMessageCB(msg model2.BarcodeQueryMessage) {
	handler.socket.WriteJSON(msg)
}

func (handler *ClientHandlerImpl) handleClientRequest(msg []byte) {

	var barcodeQueryMsg model2.BarcodeQueryMessage
	json.Unmarshal(msg, &barcodeQueryMsg)
	handler.clientBroadcast.Send(barcodeQueryMsg)
}

func (handler *ClientHandlerImpl) provideCurrentStateToClient() {
	handler.clientBroadcast.Send(model2.BarcodeQueryMessage{
		MessageType: model2.DBStateUpdateRequest,
		Payload:     db.ScannedDB,
	})
	handler.clientBroadcast.Send(model2.BarcodeQueryMessage{
		MessageType: model2.DBStateUpdateRequest,
		Payload:     db.ErrorDBRole,
	})
	handler.clientBroadcast.Send(model2.BarcodeQueryMessage{
		MessageType: model2.DBStateUpdateRequest,
		Payload:     db.DuplicatedHistoryDB,
	})
	handler.clientBroadcast.Send(model2.BarcodeQueryMessage{
		MessageType: model2.CurrentCounterUpdateRequest,
	})
	handler.clientBroadcast.Send(model2.BarcodeQueryMessage{
		MessageType: model2.TotalCounterUpdateRequest,
	})
	handler.clientBroadcast.Send(model2.BarcodeQueryMessage{
		MessageType: model2.GetNumberOfItemInListRequest,
	})
}

func (handler *ClientHandlerImpl) handle() {
	handler.provideCurrentStateToClient()
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
		handler.handleClientRequest(message)
		log.Println(mt, string(message))
	}
	defer func() {
		log.Println("Closing web socket")
		handler.socket.Close()
	}()
}
