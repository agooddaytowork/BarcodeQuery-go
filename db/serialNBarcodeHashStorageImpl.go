package db

import (
	"BarcodeQuery/classifier"
	"BarcodeQuery/model"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"log"
	"os"
	"strings"
	"time"
)

type SerialNBarcodeHashStorageImpl struct {
	DBRole              DBRole
	FilePath            string
	Store               map[string]string
	Broadcaster         *broadcast.Broadcaster
	ClientListener      *broadcast.Listener
	IgnoreClientRequest bool
}

func (db *SerialNBarcodeHashStorageImpl) GetStoreAsQueryResultArray() []QueryResult {
	var result []QueryResult

	return result
}

func (db *SerialNBarcodeHashStorageImpl) HandleClientRequest() {
	if db.ClientListener == nil {
		log.Println("Client listener is nil")
		return
	}
	for true {
		request := <-db.ClientListener.Channel()
		msg := request.(model.BarcodeQueryMessage)
		if msg.MessageType == model.DBStateUpdateRequest && len(db.Store) != 0 {
		}
	}
}

func (db *SerialNBarcodeHashStorageImpl) Load(classifier classifier.TupleClassifier) *BarcodeDBError {
	data, err := os.ReadFile(db.FilePath)

	if err != nil {
		log.Println(err)
		return nil
	}
	elements := strings.Split(string(data), "\n")
	newStorage := make(map[string]string)
	for _, e := range elements {
		key, val := classifier.Classify(e)
		if key != "" {
			newStorage[key] = val
		}
	}

	log.Printf("LOAD %d items from %s \n", len(newStorage), db.FilePath)
	db.Store = newStorage
	return nil
}

func (db *SerialNBarcodeHashStorageImpl) dump(inputPath string) *BarcodeDBError {
	//if len(db.Store) == 0 {
	//	return nil
	//}

	f, err := os.Create(inputPath)
	if err != nil {
		panic(err)
	}

	for key := range db.Store {
		_, err := f.WriteString(key + "\n")
		if err != nil {
			return &BarcodeDBError{
				ExceptionMsg: err.Error(),
			}
		}
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	return nil
}

func (db *SerialNBarcodeHashStorageImpl) Dump() *BarcodeDBError {
	return db.dump(db.FilePath)
}

func (db *SerialNBarcodeHashStorageImpl) DumpWithTimeStamp() *BarcodeDBError {
	fileName := strings.Replace(db.FilePath, ".txt", "", 1) + "-" + time.Now().Format("2006-01-02-15-04-05") + ".txt"
	return db.dump(fileName)
}

func (db *SerialNBarcodeHashStorageImpl) Insert(input string, queriedValue string) *BarcodeDBError {

	if _, ok := db.Store[input]; !ok {
		db.Store[input] = queriedValue
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
func (db *SerialNBarcodeHashStorageImpl) sendResponse(msg model.BarcodeQueryMessage) {
	if db.Broadcaster == nil {
		log.Println("DB Broadcaster is nil")
		return
	}
	db.Broadcaster.Send(msg)
}
func (db *SerialNBarcodeHashStorageImpl) Query(input string) int {
	if _, ok := db.Store[input]; ok {
		db.sendResponse(model.BarcodeQueryMessage{
			MessageType: model.DBQueryNoti,
			Payload: QueryResult{
				DBRole:      db.DBRole,
				QueryString: input,
				QueryResult: 1,
			},
		})
		return 1
	}

	db.sendResponse(model.BarcodeQueryMessage{
		MessageType: model.DBQueryNoti,
		Payload: QueryResult{
			DBRole:      db.DBRole,
			QueryString: input,
			QueryResult: -1,
		},
	})
	return -1
}

func (db *SerialNBarcodeHashStorageImpl) Clear() {
	db.Store = make(map[string]string)
}

func (db *SerialNBarcodeHashStorageImpl) GetDBLength() int {
	return len(db.Store)
}
func (db *SerialNBarcodeHashStorageImpl) GetStore() map[string]string {
	return db.Store
}

func (db *SerialNBarcodeHashStorageImpl) Sync(input map[string]string) {
}
