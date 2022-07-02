package main

import (
	"BarcodeQuery/app"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"BarcodeQuery/web"
	"github.com/textileio/go-threads/broadcast"
	"time"
)

func main() {

	dbBroadCast := broadcast.NewBroadcaster(100)
	clientBroadCast := broadcast.NewBroadcaster(100)

	existingDB := db.BarcodeDBHashStorageImpl{
		DBRole:              db.ExistingDBRole,
		FilePath:            "test/100k.txt",
		Store:               make(map[string]int),
		Broadcaster:         dbBroadCast,
		ClientListener:      clientBroadCast.Listen(),
		IgnoreClientRequest: true,
	}
	err := existingDB.Load()

	errorDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ErrorDBRole,
		FilePath:       "test/errorDB.txt",
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	duplicatedHistoryDbB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.DuplicatedHistoryDB,
		FilePath:       "test/duplicatedDB.txt",
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	scannedDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ScannedDB,
		FilePath:       "test/scannedDB.txt",
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	if err != nil {
		panic(err)
	}

	testFileReader := reader.TestFileReader{
		Interval: time.Millisecond * 200,
	}

	testFileReader.Load("test/query")

	program := app.BarcodeQueryAppImpl{
		ExistingDB:        &existingDB,
		ErrorDB:           &errorDB,
		DuplicatedItemDB:  &duplicatedHistoryDbB,
		ScannedDB:         &scannedDB,
		Reader:            &testFileReader,
		QueryCounter:      0,
		QueryCounterLimit: 100,
		Broadcaster:       dbBroadCast,
		ClientListener:    clientBroadCast.Listen(),
	}

	theWeb := web.BarcodeQueryWebImpl{
		Broadcaster:     dbBroadCast,
		ClientBroadCast: clientBroadCast,
	}

	go theWeb.Run()
	program.Run()

}
