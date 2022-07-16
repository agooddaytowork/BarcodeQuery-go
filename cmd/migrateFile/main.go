package main

import (
	"BarcodeQuery/classifier"
	"BarcodeQuery/db"
	"BarcodeQuery/model"
	"BarcodeQuery/util"
	"time"
)

func main() {
	barcodeExistingDB := db.SerialHashStorageImpl{
		DBRole:              db.ExistingDBRole,
		FilePath:            "test/persisted-backup.txt",
		Store:               make(map[string]int),
		Broadcaster:         nil,
		ClientListener:      nil,
		IgnoreClientRequest: true,
	}

	barcodeExistingDB.Load(&classifier.DummyBarcodeTupleClassifier{})

	newPersisted := make(map[string]model.PersistedSerialRecord)

	for k := range barcodeExistingDB.Store {
		newPersisted[k] = model.PersistedSerialRecord{
			Serial:           k,
			ScannedTimestamp: time.Now().Unix(),
			Lot:              "Từ file cũ",
		}
	}
	util.DumpConfigToFile("test/persisted-new.txt", newPersisted)
}
