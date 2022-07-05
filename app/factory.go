package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"time"
)

func GetBarcodeQueryAppImpl(theConfig BarcodeAppConfig, dbBroadCast *broadcast.Broadcaster, clientBroadCast *broadcast.Broadcaster) BarcodeQueryAppImpl {
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
	case TestFileReader:
		testFileReader := reader.TestFileReader{
			Interval: time.Millisecond * 200,
		}
		testFileReader.Load(theConfig.ReaderURI)
		theReader = &testFileReader

	case ConsoleReader:
		theReader = &reader.ConsoleReader{}

	case TCPReader:
		theReader = &reader.TCPReader{
			URL:           theConfig.ReaderURI,
			SpawnedThread: false,
			ReportChannel: make(chan string, 1000),
		}
	default:
		panic(fmt.Sprintf("Unsupported reader, only support %s/%s/%s", TestFileReader, ConsoleReader, TCPReader))
	}

	// init the program
	return BarcodeQueryAppImpl{
		ExistingDB:        &existingDB,
		ErrorDB:           &errorDB,
		DuplicatedItemDB:  &duplicatedHistoryDbB,
		ScannedDB:         &scannedDB,
		Reader:            theReader,
		QueryCounter:      0,
		QueryCounterLimit: theConfig.QueryCounterLimit,
		TotalCounter:      0,
		Broadcaster:       dbBroadCast,
		ClientListener:    clientBroadCast.Listen(),
		Actuator:          &actuator.ConsoleActuator{},
	}

}
