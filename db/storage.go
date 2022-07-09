package db

type BarcodeDB interface {
	Load() *BarcodeDBError
	Dump() *BarcodeDBError
	DumpWithTimeStamp() *BarcodeDBError
	Insert(input string, queriedValue int) *BarcodeDBError
	Query(input string) int
	Clear()
	HandleClientRequest()
	GetStoreAsQueryResultArray() []QueryResult
	GetDBLength() int
}
