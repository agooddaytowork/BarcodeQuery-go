package web

import (
	"BarcodeQuery/app"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type BarcodeQueryWebImpl struct {
	handler ClientHandler
}

func (web *BarcodeQueryWebImpl) Run() {

}

var upgrade = websocket.Upgrader{}

func (web *BarcodeQueryWebImpl) barcodeWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go web.handler.handle(c)
}

func (web *BarcodeQueryWebImpl) RegisterDBCallBack(dbRole app.DBRole, callback func()) {
	http.HandleFunc("/ws", web.barcodeWS)
}
