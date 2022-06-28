package app

import (
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type BarcodeQueryAppImpl struct {
	ExistingDB        db.BarcodeDB
	QueriedHistoryDB  db.BarcodeDB
	ErrorDB           db.BarcodeDB
	Reader            reader.BarcodeReader
	QueryCounter      int
	QueryCounterLimit int
}

func (app *BarcodeQueryAppImpl) Run() {

	run := true
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		print("common")
		run = false
	}()

	for run {
		queryString := app.Reader.Read()
		queryResult := app.ExistingDB.Query(queryString)

		if queryResult < 0 {
			// not found in existing DB

			errorQuery := app.ErrorDB.Query(queryString)

			if errorQuery == -1 {
				app.ErrorDB.Insert(queryString, 0)
			}

		} else if queryResult == 1 {
			// found barcode
			// do something
			app.QueryCounter++
		} else {
			// found duplicated query
			duplicateQuery := app.QueriedHistoryDB.Query(queryString)

			if duplicateQuery == -1 {
				app.QueriedHistoryDB.Insert(queryString, 0)
			}

		}

		if app.QueryCounter == app.QueryCounterLimit {
			app.ErrorDB.DumpWithTimeStamp()
			app.QueriedHistoryDB.DumpWithTimeStamp()
		}

		fmt.Printf("Query result %s : %d \n", queryString, queryResult)
	}

	defer func() {
		// Todo: dump these when DBCounter hit limit as well
		fmt.Println("Cleaning up")
		app.ErrorDB.DumpWithTimeStamp()
		app.QueriedHistoryDB.DumpWithTimeStamp()
	}()
}
