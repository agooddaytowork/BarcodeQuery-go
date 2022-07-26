package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/classifier"
	"BarcodeQuery/db"
	"BarcodeQuery/hashing"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"time"
)

func GetBarcodeQueryAppImpl(configPath string, theConfig BarcodeAppConfig, dbBroadCast *broadcast.Broadcaster, clientBroadCast *broadcast.Broadcaster, config BarcodeAppConfig) BarcodeQueryAppImpl {

	persistedScanDB := db.PersistedSerialRecordHashStorageImpl{
		DBRole:              db.PersitedDBRole,
		FilePath:            "persisted.txt",
		Store:               make(map[string]model.PersistedSerialRecord),
		Broadcaster:         nil,
		ClientListener:      nil,
		IgnoreClientRequest: true,
	}
	err := persistedScanDB.Load(&classifier.DummyBarcodeTupleClassifier{})
	if err != nil {
		panic(err)
	}

	barcodeExistingDB := db.SerialHashStorageImpl{
		DBRole:              db.ExistingDBRole,
		FilePath:            theConfig.ExistingDBPath,
		Store:               make(map[string]int),
		Broadcaster:         dbBroadCast,
		ClientListener:      clientBroadCast.Listen(),
		IgnoreClientRequest: true,
	}
	err = barcodeExistingDB.Load(&classifier.BarcodeTupleClassifier{})
	if err != nil {
		panic(err)
	}

	barcodeNSerialDB := db.SerialNBarcodeHashStorageImpl{
		DBRole:              db.BarcodeVsSerialDB,
		FilePath:            theConfig.ExistingDBPath,
		Store:               make(map[string]string),
		Broadcaster:         dbBroadCast,
		ClientListener:      clientBroadCast.Listen(),
		IgnoreClientRequest: true,
	}
	barcodeNSerialDB.Load(&classifier.BarcodeNSerialTupleClassifier{})

	serialNBarcodeDB := db.SerialNBarcodeHashStorageImpl{
		DBRole:              db.BarcodeVsSerialDB,
		FilePath:            theConfig.ExistingDBPath,
		Store:               make(map[string]string),
		Broadcaster:         dbBroadCast,
		ClientListener:      clientBroadCast.Listen(),
		IgnoreClientRequest: true,
	}

	serialNBarcodeDB.Load(&classifier.SerialNBarcodeTupleClassifier{})

	errorDB := db.SerialHashStorageImpl{
		DBRole:         db.ErrorDBRole,
		FilePath:       theConfig.ErrorDBPath,
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	duplicatedHistoryDbB := db.SerialHashStorageImpl{
		DBRole:         db.DuplicatedHistoryDB,
		FilePath:       theConfig.DuplicatedDBPath,
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	scannedDB := db.SerialHashStorageImpl{
		DBRole:         db.ScannedDB,
		FilePath:       theConfig.ScannedDBPath,
		Store:          make(map[string]int),
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
	}

	actuator := actuator.GetActuator(theConfig.ActuatorType, theConfig.ActuatorURI)

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
			URL:                         theConfig.ReaderURI,
			SpawnedThread:               false,
			ReportChannel:               make(chan string, 1000),
			DuplicateDebounceIntervalMs: theConfig.ReaderDuplicateDebounceIntervalMs,
		}
	default:
		panic(fmt.Sprintf("Unsupported reader, only support %s/%s/%s", reader.TestFileReaderType, reader.ConsoleReaderType, reader.TCPReaderType))
	}
	// init the program
	return BarcodeQueryAppImpl{
		PersistedScannedDB: &persistedScanDB,
		BarcodeExistingDB:  &barcodeExistingDB,
		ErrorDB:            &errorDB,
		DuplicatedItemDB:   &duplicatedHistoryDbB,
		BarcodeAndSerialDB: &barcodeNSerialDB,
		SerialAndBarcodeDB: &serialNBarcodeDB,
		ScannedDB:          &scannedDB,
		Reader:             theReader,
		ConfigPath:         configPath,
		CounterReport: model.CounterReport{
			QueryCounter:             0,
			QueryCounterLimit:        theConfig.QueryCounterLimit,
			TotalCounter:             0,
			PackageCounter:           0,
			NumberOfItemInExistingDB: 0,
			NumberOfCameraScanError:  0,
		},
		Broadcaster:    dbBroadCast,
		ClientListener: clientBroadCast.Listen(),
		Actuator:       actuator,
		Config:         config,
		Hasher:         &hashing.BarcodeSHA256HasherImpl{},
		TestMode:       false,
	}

}
