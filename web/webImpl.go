package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/textileio/go-threads/broadcast"
	"io"
	"log"
	"net/http"
	"os"
)

type BarcodeQueryWebImpl struct {
	Broadcaster         *broadcast.Broadcaster
	ClientBroadCast     *broadcast.Broadcaster
	StaticFilePath      string
	BarcodeListFilePath string
}

func (web *BarcodeQueryWebImpl) Run() {
	fs := http.FileServer(http.Dir(web.StaticFilePath))
	http.Handle("/", fs)
	http.HandleFunc("/ws", web.barcodeWS)
	http.HandleFunc("/uploadList", web.handleUploadList)
	http.ListenAndServe(":80", nil)
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (web *BarcodeQueryWebImpl) handleUploadList(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("danhsach")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(web.BarcodeListFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(200)
	}
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
