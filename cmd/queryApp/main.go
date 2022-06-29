package main

import (
	"BarcodeQuery/app"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"BarcodeQuery/web"
	"github.com/textileio/go-threads/broadcast"
)

func main() {

	msgBroadCast := broadcast.NewBroadcaster(100)

	existingDB := db.BarcodeDBHashStorageImpl{
		DBRole:      db.ExistingDBRole,
		FilePath:    "/Users/tamduong/Workspace/duc/BarcodeQuery-go/test/100k.txt",
		Store:       make(map[string]int),
		Broadcaster: msgBroadCast,
	}
	err := existingDB.Load()

	errorDB := db.BarcodeDBHashStorageImpl{
		DBRole:      db.ErrorDBRole,
		FilePath:    "/Users/tamduong/Workspace/duc/BarcodeQuery-go/test/blabla.txt",
		Store:       make(map[string]int),
		Broadcaster: msgBroadCast,
	}

	queriedHistoryDB := db.BarcodeDBHashStorageImpl{
		DBRole:      db.QueriedHistoryDBRole,
		FilePath:    "/Users/tamduong/Workspace/duc/BarcodeQuery-go/test/bloblo.txt",
		Store:       make(map[string]int),
		Broadcaster: msgBroadCast,
	}

	if err != nil {
		panic(err)
	}
	program := app.BarcodeQueryAppImpl{
		ExistingDB:        &existingDB,
		ErrorDB:           &errorDB,
		QueriedHistoryDB:  &queriedHistoryDB,
		Reader:            &reader.ConsoleReader{},
		QueryCounter:      0,
		QueryCounterLimit: 10,
		Broadcaster:       msgBroadCast,
	}

	go program.Run()

	theWeb := web.BarcodeQueryWebImpl{
		Broadcaster: msgBroadCast,
	}

	theWeb.Run()

}
