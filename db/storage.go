package db

import (
	"BarcodeQuery/classifier"
	"BarcodeQuery/model"
)

type SerialDB interface {
	Load(classifier classifier.TupleClassifier) *BarcodeDBError
	Dump() *BarcodeDBError
	DumpWithTimeStamp() *BarcodeDBError
	Insert(input string, queriedValue int) *BarcodeDBError
	Query(input string) int
	Clear()
	HandleClientRequest()
	GetStoreAsQueryResultArray() []QueryIntResult
	GetDBLength() int
	Sync(input map[string]int)
	GetStore() map[string]int
}

type SerialNBarcodeDB interface {
	Load(classifier classifier.TupleClassifier) *BarcodeDBError
	Dump() *BarcodeDBError
	DumpWithTimeStamp() *BarcodeDBError
	Insert(input string, queriedValue string) *BarcodeDBError
	Query(input string) string
	Clear()
	HandleClientRequest()
	GetStoreAsQueryResultArray() []QueryIntResult
	GetDBLength() int
	Sync(input map[string]string)
	GetStore() map[string]string
}

type PersistedSerialRecordDB interface {
	Load(classifier classifier.TupleClassifier) *BarcodeDBError
	Dump() *BarcodeDBError
	DumpWithTimeStamp() *BarcodeDBError
	Insert(input string, v model.PersistedSerialRecord) *BarcodeDBError
	Query(input string) (model.PersistedSerialRecord, bool)
	Clear()
	HandleClientRequest()
	GetStoreAsQueryResultArray() []QueryIntResult
	GetDBLength() int
	Sync(input map[string]model.PersistedSerialRecord)
	GetStore() map[string]model.PersistedSerialRecord
}
