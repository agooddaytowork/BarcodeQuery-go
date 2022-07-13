package db

import "BarcodeQuery/classifier"

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
