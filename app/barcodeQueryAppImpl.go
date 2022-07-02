package app

import (
	"BarcodeQuery/db"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"log"
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
	Broadcaster       *broadcast.Broadcaster
	ClientListener    *broadcast.Listener
}

func (app *BarcodeQueryAppImpl) handleClientRequest() {
	for true {
		request := <-app.ClientListener.Channel()

		msg := request.(model.BarcodeQueryMessage)

		if msg.MessageType == model.CurrentCounterUpdateRequest {
			app.Broadcaster.Send(
				model.BarcodeQueryMessage{
					MessageType: model.CurrentCounterUpdateResponse,
					Payload:     app.QueryCounter,
				})
		}
	}
}

func (app *BarcodeQueryAppImpl) cleanUp() {
	log.Println("Cleaning up")
	app.QueryCounter = 0
	app.ScannedDB.DumpWithTimeStamp()

	app.ErrorDB.DumpWithTimeStamp()
	app.DuplicatedItemDB.DumpWithTimeStamp()

	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
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

	go app.handleClientRequest()
	go app.DuplicatedItemDB.HandleClientRequest()
	go app.ErrorDB.HandleClientRequest()
	go app.ScannedDB.HandleClientRequest()

	for run {
		queryString := app.Reader.Read()
		existingDBResult := app.ExistingDB.Query(queryString)

		if existingDBResult < 0 {
			// not found in existing DB

			errorQuery := app.ErrorDB.Query(queryString)

			if errorQuery == -1 {
				app.ErrorDB.Insert(queryString, 0)
			}

		} else if existingDBResult == 1 {
			// found barcode
			// do something
			app.ScannedDB.Insert(queryString, 0)
			app.ScannedDB.Query(queryString)
			app.QueryCounter++
		} else {
			// found duplicated query
			duplicateQuery := app.DuplicatedItemDB.Query(queryString)

			if duplicateQuery == -1 {
				app.DuplicatedItemDB.Insert(queryString, 0)
			}
		}

		if app.QueryCounter == app.QueryCounterLimit {
			app.cleanUp()
		}

		app.Broadcaster.Send(model.BarcodeQueryMessage{
			MessageType: model.CounterNoti,
			Payload:     app.QueryCounter,
		})

		fmt.Printf("Query result %s : %d \n", queryString, existingDBResult)
	}

	defer app.cleanUp()
}
