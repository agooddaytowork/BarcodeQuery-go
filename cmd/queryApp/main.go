package main

import (
	"BarcodeQuery/app"
	"BarcodeQuery/config"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"BarcodeQuery/web"
	"flag"
	"github.com/textileio/go-threads/broadcast"
	"time"
)

func main() {

	configPath := flag.String("c", "test/config.json", "Config path")
	flag.Parse()

	theConfig := config.LoadConfigFromFile(*configPath)

	dbBroadCast := broadcast.NewBroadcaster(100)
	clientBroadCast := broadcast.NewBroadcaster(100)

	existingDB := db.BarcodeDBHashStorageImpl{
		DBRole:              db.ExistingDBRole,
		FilePath:            theConfig.ExistingDBPath,
		Store:               make(map[string]int),
		Broadcaster:         dbBroadCast,
		ClientListener:      clientBroadCast.Listen(),
		IgnoreClientRequest: true,
	}
	err := existingDB.Load()

	errorDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ErrorDBRole,
		FilePath:       theConfig.ErrorDBPath,
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	duplicatedHistoryDbB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.DuplicatedHistoryDB,
		FilePath:       theConfig.DuplicatedDBPath,
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	scannedDB := db.BarcodeDBHashStorageImpl{
		DBRole:         db.ScannedDB,
		FilePath:       theConfig.ScannedDBPath,
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	if err != nil {
		panic(err)
	}

	// get the reader
	var theReader reader.BarcodeReader
	switch theConfig.ReaderType {
	case "test_file":
		testFileReader := reader.TestFileReader{
			Interval: time.Millisecond * 200,
		}
		testFileReader.Load(theConfig.ReaderURI)
		theReader = &testFileReader
	default:
		panic("Unsupported reader")
	}

	// init the program
	program := app.BarcodeQueryAppImpl{
		ExistingDB:        &existingDB,
		ErrorDB:           &errorDB,
		DuplicatedItemDB:  &duplicatedHistoryDbB,
		ScannedDB:         &scannedDB,
		Reader:            theReader,
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
