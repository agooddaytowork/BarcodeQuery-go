package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/db"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"time"
)

func GetBarcodeQueryAppImpl(configPath string, theConfig BarcodeAppConfig, dbBroadCast *broadcast.Broadcaster, clientBroadCast *broadcast.Broadcaster, config BarcodeAppConfig) BarcodeQueryAppImpl {

	persistedScanDB := db.BarcodeDBHashStorageImpl{
		DBRole:              db.ExistingDBRole,
		FilePath:            "persisted.txt",
		Store:               make(map[string]int),
		Broadcaster:         nil,
		ClientListener:      nil,
		IgnoreClientRequest: true,
	}
	persistedScanDB.Load()

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
	case reader.TestFileReaderType:
		testFileReader := reader.TestFileReader{
			Interval: time.Millisecond * 200,
		}
		testFileReader.Load(theConfig.ReaderURI)
		theReader = &testFileReader

	case reader.ConsoleReaderType:
		theReader = &reader.ConsoleReader{}

	case reader.TCPReaderType:
		theReader = &reader.TCPReader{
			URL:           theConfig.ReaderURI,
			SpawnedThread: false,
			ReportChannel: make(chan string, 1000),
		}
	default:
		panic(fmt.Sprintf("Unsupported reader, only support %s/%s/%s", reader.TestFileReaderType, reader.ConsoleReaderType, reader.TCPReaderType))
	}
	// init the program
	return BarcodeQueryAppImpl{
		PersistedScannedDB: &persistedScanDB,
		ExistingDB:         &existingDB,
		ErrorDB:            &errorDB,
		DuplicatedItemDB:   &duplicatedHistoryDbB,
		ScannedDB:          &scannedDB,
		Reader:             theReader,
		ConfigPath:         configPath,
		CounterReport: model.CounterReport{
			QueryCounter:             0,
			QueryCounterLimit:        theConfig.QueryCounterLimit,
			TotalCounter:             0,
			PackageCounter:           0,
			NumberOfItemInExistingDB: 0,
		},
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
		Actuator:       &actuator.ConsoleActuator{},
		Config:         config,
	}

}
