package db

import (
	"BarcodeQuery/classifier"
	"BarcodeQuery/model"
	"BarcodeQuery/util"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"log"
	"strings"
	"time"
)

type PersistedSerialRecordHashStorageImpl struct {
	DBRole              DBRole
	FilePath            string
	Store               map[string]model.PersistedSerialRecord
	Broadcaster         *broadcast.Broadcaster
	ClientListener      *broadcast.Listener
	IgnoreClientRequest bool
}

func (db *PersistedSerialRecordHashStorageImpl) GetStoreAsQueryResultArray() []QueryIntResult {
	var result []QueryIntResult
	for _, v := range db.Store {
		result = append(result, QueryIntResult{
			DBRole:      db.DBRole,
			QueryString: v.Serial,
			QueryResult: 1,
		})
	}
	return result
}

func (db *PersistedSerialRecordHashStorageImpl) HandleClientRequest() {
	if db.ClientListener == nil {
		log.Println("Client listener is nil")
		return
	}
	for true {
		request := <-db.ClientListener.Channel()
		msg := request.(model.BarcodeQueryMessage)
		if msg.MessageType == model.DBStateUpdateRequest && len(db.Store) != 0 {
			if msg.Payload.(DBRole) == db.DBRole {
				db.Broadcaster.Send(
					model.BarcodeQueryMessage{
						MessageType: model.DBStateUpdateResponse,
						Payload: StateUpdate{
							DBRole: db.DBRole,
							State:  db.GetStoreAsQueryResultArray(),
						},
					},
				)
			}
		}
	}
}

func (db *PersistedSerialRecordHashStorageImpl) Load(classifier classifier.TupleClassifier) *BarcodeDBError {
	var newData map[string]model.PersistedSerialRecord
	util.LoadConfigFromFile(db.FilePath, &newData)
	log.Printf("LOAD %d items from %s \n", len(newData), db.FilePath)
	db.Store = newData
	return nil
}

func (db *PersistedSerialRecordHashStorageImpl) dump(inputPath string) *BarcodeDBError {
	util.DumpConfigToFile(inputPath, db.Store)
	return nil
}

func (db *PersistedSerialRecordHashStorageImpl) Dump() *BarcodeDBError {
	return db.dump(db.FilePath)
}

func (db *PersistedSerialRecordHashStorageImpl) DumpWithTimeStamp() *BarcodeDBError {
	fileName := strings.Replace(db.FilePath, ".txt", "", 1) + "-" + time.Now().Format("2006-01-02-15-04-05") + ".txt"
	return db.dump(fileName)
}

func (db *PersistedSerialRecordHashStorageImpl) Insert(input string, v model.PersistedSerialRecord) *BarcodeDBError {
	if _, ok := db.Store[input]; !ok {
		db.Store[input] = v
		return nil
	}
	return &BarcodeDBError{
		ExceptionMsg: fmt.Sprintf("value %s exist", input),
	}
}

/*
Query
Return -1 if reader not found in list
Return the number this reader has been queried
*/
func (db *PersistedSerialRecordHashStorageImpl) sendResponse(msg model.BarcodeQueryMessage) {
	if db.Broadcaster == nil {
		log.Println("DB Broadcaster is nil")
		return
	}
	db.Broadcaster.Send(msg)
}
func (db *PersistedSerialRecordHashStorageImpl) Query(input string) (model.PersistedSerialRecord, bool) {
	if v, ok := db.Store[input]; ok {
		db.sendResponse(model.BarcodeQueryMessage{
			MessageType: model.DBQueryNoti,
			Payload: QueryIntResult{
				DBRole:      db.DBRole,
				QueryString: input,
				QueryResult: 1,
			},
		})
		return v, ok
	}

	db.sendResponse(model.BarcodeQueryMessage{
		MessageType: model.DBQueryNoti,
		Payload: QueryIntResult{
			DBRole:      db.DBRole,
			QueryString: input,
			QueryResult: -1,
		},
	})
	return model.PersistedSerialRecord{}, false
}

func (db *PersistedSerialRecordHashStorageImpl) Clear() {
	db.Store = make(map[string]model.PersistedSerialRecord)
}

func (db *PersistedSerialRecordHashStorageImpl) GetDBLength() int {
	return len(db.Store)
}
func (db *PersistedSerialRecordHashStorageImpl) GetStore() map[string]model.PersistedSerialRecord {
	return db.Store
}

func (db *PersistedSerialRecordHashStorageImpl) Sync(input map[string]model.PersistedSerialRecord) {
	log.Println("PersistedSerialRecordHashStorageImpl Sync not implemented yet")
}
