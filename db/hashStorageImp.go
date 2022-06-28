package db

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type BarcodeDBHashStorageImpl struct {
	FilePath        string
	Store           map[string]int
	DBQueryCallBack func(query string, queryCounter int)
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

func (db *BarcodeDBHashStorageImpl) dump(inputPath string) *BarcodeDBError {
	if len(db.Store) == 0 {
		return nil
	}

	f, err := os.Create(inputPath)
	if err != nil {
		return &BarcodeDBError{
			ExceptionMsg: err.Error(),
		}
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
func (db *BarcodeDBHashStorageImpl) Query(input string) int {
	if queriedNumber, ok := db.Store[input]; ok {
		newQueriedNumber := queriedNumber + 1
		db.Store[input] = newQueriedNumber

		go db.DBQueryCallBack(input, newQueriedNumber)
		return newQueriedNumber
	}

	go db.DBQueryCallBack(input, -1)
	return -1
}
