package db

type BarcodeDB interface {
	Load() *BarcodeDBError
	Dump() *BarcodeDBError
	Insert(input string, queriedValue int) *BarcodeDBError
	Query(input string) int
}
