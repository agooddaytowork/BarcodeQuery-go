package db

type BarcodeDB interface {
	Load() *BarcodeDBError
	Dump() *BarcodeDBError
	Insert(input string) *BarcodeDBError
	Query(input string) int
}
