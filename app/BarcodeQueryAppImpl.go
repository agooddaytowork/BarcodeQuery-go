package app

import (
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"fmt"
)

type BarcodeQueryAppImpl struct {
	Db     db.BarcodeDB
	Reader reader.BarcodeReader
}

func (app *BarcodeQueryAppImpl) Run() {

	for true {
		queryString := app.Reader.Read()
		queryResult := app.Db.Query(queryString)
		fmt.Printf("Query result %s : %d \n", queryString, queryResult)
	}
}
