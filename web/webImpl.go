package web

import (
	"github.com/gorilla/websocket"
	"github.com/textileio/go-threads/broadcast"
	"log"
	"net/http"
)

type BarcodeQueryWebImpl struct {
	Broadcaster     *broadcast.Broadcaster
	ClientBroadCast *broadcast.Broadcaster
	StaticFilePath  string
}

func (web *BarcodeQueryWebImpl) Run() {
	fs := http.FileServer(http.Dir(web.StaticFilePath))
	http.Handle("/", fs)
	http.HandleFunc("/ws", web.barcodeWS)
	http.ListenAndServe(":80", nil)
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (web *BarcodeQueryWebImpl) barcodeWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	handler := ClientHandlerImpl{
		socket:          c,
		dbListener:      web.Broadcaster.Listen(),
		clientBroadcast: web.ClientBroadCast,
	}

	go handler.handle()
}
