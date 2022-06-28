package db

import (
	"fmt"
	"os"
	"strings"
)

type BarcodeDBHashStorageImpl struct {
	FilePath string
	Store    map[string]int
}

func (db *BarcodeDBHashStorageImpl) Load() *BarcodeDBError {
	data, err := os.ReadFile(db.FilePath)

	if err != nil {

		return &BarcodeDBError{
			ExceptionMsg: err.Error(),
		}
	}

	elements := strings.Split(string(data), "\n")

	newStorage := make(map[string]int)
	for _, e := range elements {
		newStorage[strings.Trim(e, " ")] = 0
	}

	db.Store = newStorage
	return nil
}

func (db *BarcodeDBHashStorageImpl) Dump() *BarcodeDBError {
	f, err := os.Create(db.FilePath)

	if err != nil {

		return &BarcodeDBError{
			ExceptionMsg: err.Error(),
		}
	}

	for key, element := range db.Store {
		fmt.Println("Key:", key, "=>", "Element:", element)
		_, err := f.WriteString(key)
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

func (db *BarcodeDBHashStorageImpl) Insert(input string) *BarcodeDBError {

	if _, ok := db.Store[input]; !ok {
		db.Store[input] = 0
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
func (db *BarcodeDBHashStorageImpl) Query(input string) int {
	if queriedNumber, ok := db.Store[input]; ok {
		newQueriedNumber := queriedNumber + 1
		db.Store[input] = newQueriedNumber

		return newQueriedNumber
	}
	return -1
}
