package app

import (
	"BarcodeQuery/db"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"os"
	"os/signal"
	"syscall"
)

type BarcodeQueryAppImpl struct {
	ExistingDB        db.BarcodeDB
	DuplicatedItemDB  db.BarcodeDB
	ErrorDB           db.BarcodeDB
	ScannedDB         db.BarcodeDB
	Reader            reader.BarcodeReader
	QueryCounter      int
	QueryCounterLimit int
	DBBroadcast       *broadcast.Broadcaster
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

	go app.ExistingDB.HandleClientRequest()
	go app.DuplicatedItemDB.HandleClientRequest()
	go app.ErrorDB.HandleClientRequest()

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
			app.ScannedDB.Insert(queryString, 0)
			app.QueryCounter++
		} else {
			// found duplicated query
			duplicateQuery := app.DuplicatedItemDB.Query(queryString)

			if duplicateQuery == -1 {
				app.DuplicatedItemDB.Insert(queryString, 0)
			}
		}

		if app.QueryCounter == app.QueryCounterLimit {
			app.QueryCounter = 0
			app.ScannedDB.DumpWithTimeStamp()

			app.ErrorDB.DumpWithTimeStamp()
			app.DuplicatedItemDB.DumpWithTimeStamp()

			app.ScannedDB.Clear()
			app.ErrorDB.Clear()
			app.DuplicatedItemDB.Clear()
		}

		app.DBBroadcast.Send(model.BarcodeQueryMessage{
			MessageType: model.CounterNoti,
			Payload:     app.QueryCounter,
		})

		fmt.Printf("Query result %s : %d \n", queryString, queryResult)
	}

	defer func() {
		// Todo: dump these when DBCounter hit limit as well
		fmt.Println("Cleaning up")
		app.ScannedDB.DumpWithTimeStamp()
		app.ErrorDB.DumpWithTimeStamp()
		app.DuplicatedItemDB.DumpWithTimeStamp()
	}()
}
