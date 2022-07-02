package main

import (
	"BarcodeQuery/app"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"BarcodeQuery/web"
	"github.com/textileio/go-threads/broadcast"
)

func main() {

	dbBroadCast := broadcast.NewBroadcaster(100)
	clientBroadCast := broadcast.NewBroadcaster(100)

	existingDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ExistingDBRole,
		FilePath:       "test/100k.txt",
		Store:          make(map[string]int),
		DBBroadCast:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}
	err := existingDB.Load()

	errorDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ErrorDBRole,
		FilePath:       "test/errorDB.txt",
		Store:          make(map[string]int),
		DBBroadCast:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	duplicatedHistoryDbB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.DuplicatedHistoryDB,
		FilePath:       "test/duplicatedDB.txt",
		Store:          make(map[string]int),
		DBBroadCast:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	scannedDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ScannedDB,
		FilePath:       "test/scannedDB.txt",
		Store:          make(map[string]int),
		DBBroadCast:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	if err != nil {
		panic(err)
	}
	program := app.BarcodeQueryAppImpl{
		ExistingDB:        &existingDB,
		ErrorDB:           &errorDB,
		DuplicatedItemDB:  &duplicatedHistoryDbB,
		ScannedDB:         &scannedDB,
		Reader:            &reader.ConsoleReader{},
		QueryCounter:      0,
		QueryCounterLimit: 10,
		DBBroadcast:       dbBroadCast,
	}

	go program.Run()

	theWeb := web.BarcodeQueryWebImpl{
		DBBroadcast:     dbBroadCast,
		ClientBroadCast: clientBroadCast,
	}

	theWeb.Run()

}
