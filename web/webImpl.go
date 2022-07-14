package web

import (
	"BarcodeQuery/classifier"
	"BarcodeQuery/db"
	"BarcodeQuery/hashing"
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

func hashingSerialNBarcodeProductionFile(filePath string) {
	log.Println("hashingSerialNBarcodeProductionFile")
	barcodeNSerialDB := db.SerialNBarcodeHashStorageImpl{
		DBRole:              db.BarcodeVsSerialDB,
		FilePath:            filePath,
		Store:               make(map[string]string),
		Broadcaster:         nil,
		ClientListener:      nil,
		IgnoreClientRequest: true,
	}

	barcodeNSerialDB.Load(&classifier.SerialNBarcodeTupleClassifier{})
	hasher := hashing.BarcodeSHA256HasherImpl{}
	newStore := make(map[string]string)

	for k, e := range barcodeNSerialDB.GetStore() {
		hashValue := hasher.Hash(e)
		log.Printf("%s : %s", e, hashValue)
		newStore[k] = hashValue
	}
	barcodeNSerialDB.Sync(newStore)
	barcodeNSerialDB.Dump()
	log.Println("hashing complete")

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (web *BarcodeQueryWebImpl) handleUploadList(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("handleUploadList")
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
		hashingSerialNBarcodeProductionFile(web.BarcodeListFilePath)
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
