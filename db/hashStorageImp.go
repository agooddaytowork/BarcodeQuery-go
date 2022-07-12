package db

import (
	"BarcodeQuery/model"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"log"
	"os"
	"strings"
	"time"
)

type BarcodeDBHashStorageImpl struct {
	DBRole              DBRole
	FilePath            string
	Store               map[string]int
	Broadcaster         *broadcast.Broadcaster
	ClientListener      *broadcast.Listener
	IgnoreClientRequest bool
}

func (db *BarcodeDBHashStorageImpl) GetStoreAsQueryResultArray() []QueryResult {
	var result []QueryResult

	for element := range db.Store {
		result = append(result, QueryResult{
			DBRole:      db.DBRole,
			QueryString: element,
			QueryResult: db.Store[element],
		})
	}

	return result
}

func (db *BarcodeDBHashStorageImpl) HandleClientRequest() {
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

func (db *BarcodeDBHashStorageImpl) Load() *BarcodeDBError {
	data, err := os.ReadFile(db.FilePath)

	if err != nil {
		log.Println(err)
		return nil
	}
	elements := strings.Split(string(data), "\n")
	newStorage := make(map[string]int)
	for _, e := range elements {
		element := strings.Trim(e, " \r\t")
		if element != "" {
			newStorage[element] = 0
		}
	}

	log.Printf("LOAD %d items from %s \n", len(newStorage), db.FilePath)
	db.Store = newStorage
	return nil
}

func (db *BarcodeDBHashStorageImpl) dump(inputPath string) *BarcodeDBError {
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

func (db *BarcodeDBHashStorageImpl) Dump() *BarcodeDBError {
	return db.dump(db.FilePath)
}

func (db *BarcodeDBHashStorageImpl) DumpWithTimeStamp() *BarcodeDBError {
	fileName := strings.Replace(db.FilePath, ".txt", "", 1) + "-" + time.Now().Format("2006-01-02-15-04-05") + ".txt"
	return db.dump(fileName)
}

func (db *BarcodeDBHashStorageImpl) Insert(input string, queriedValue int) *BarcodeDBError {

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
func (db *BarcodeDBHashStorageImpl) sendResponse(msg model.BarcodeQueryMessage) {
	if db.Broadcaster == nil {
		log.Println("DB Broadcaster is nil")
		return
	}
	db.Broadcaster.Send(msg)
}
func (db *BarcodeDBHashStorageImpl) Query(input string) int {
	if queriedNumber, ok := db.Store[input]; ok {
		newQueriedNumber := queriedNumber + 1
		db.Store[input] = newQueriedNumber
		db.sendResponse(model.BarcodeQueryMessage{
			MessageType: model.DBQueryNoti,
			Payload: QueryResult{
				DBRole:      db.DBRole,
				QueryString: input,
				QueryResult: newQueriedNumber,
			},
		})
		return newQueriedNumber
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

func (db *BarcodeDBHashStorageImpl) Clear() {
	db.Store = make(map[string]int)
}

func (db *BarcodeDBHashStorageImpl) GetDBLength() int {
	return len(db.Store)
}
func (db *BarcodeDBHashStorageImpl) GetStore() map[string]int {
	return db.Store
}

func (db *BarcodeDBHashStorageImpl) Sync(input map[string]int) {
	for key := range input {
		db.Store[key] = 1
	}
}
