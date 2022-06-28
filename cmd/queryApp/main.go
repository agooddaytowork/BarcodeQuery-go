package main

import (
	app2 "BarcodeQuery/app"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"fmt"
)

func callBack(input string, value int) {
	fmt.Println("callBack", input, " ", value)
}

func main() {
	existingDB := db.BarcodeDBHashStorageImpl{
		FilePath:        "/Users/tam/Workspace/Duc/BarcodeQuery/test/100k.txt",
		Store:           make(map[string]int),
		DBQueryCallBack: callBack,
	}
	err := existingDB.Load()

	errorDB := db.BarcodeDBHashStorageImpl{
		FilePath:        "/Users/tam/Workspace/Duc/BarcodeQuery/test/errorDB.txt",
		Store:           make(map[string]int),
		DBQueryCallBack: callBack,
	}

	queriedHistoryDB := db.BarcodeDBHashStorageImpl{
		FilePath:        "/Users/tam/Workspace/Duc/BarcodeQuery/test/queriedHistoryDB.txt",
		Store:           make(map[string]int),
		DBQueryCallBack: callBack,
	}

	if err != nil {
		panic(err)
	}
	program := app2.BarcodeQueryAppImpl{
		ExistingDB:        &existingDB,
		ErrorDB:           &errorDB,
		QueriedHistoryDB:  &queriedHistoryDB,
		Reader:            &reader.ConsoleReader{},
		QueryCounter:      0,
		QueryCounterLimit: 10,
	}

	program.Run()

}
