package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/db"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
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
	TotalCounter      int
	Broadcaster       *broadcast.Broadcaster
	ClientListener    *broadcast.Listener
	Actuator          actuator.BarcodeActuator
}

func (app *BarcodeQueryAppImpl) sendResponse(msgType model.MessageType, payload any) {
	app.Broadcaster.Send(
		model.BarcodeQueryMessage{
			MessageType: msgType,
			Payload:     payload,
		})
}

func (app *BarcodeQueryAppImpl) handleClientRequest() {
	for true {
		request := <-app.ClientListener.Channel()

		msg := request.(model.BarcodeQueryMessage)

		switch msg.MessageType {

		case model.CurrentCounterUpdateRequest:
			app.sendResponse(model.CurrentCounterUpdateResponse, app.QueryCounter)
		case model.TotalCounterUpdateRequest:
			app.sendResponse(model.TotalCounterUpdateResponse, app.TotalCounter)
		case model.SetErrorActuatorRequest:
			state := msg.Payload.(actuator.ActuatorState)
			app.Actuator.SetErrorActuatorState(state)
			app.sendResponse(model.SetErrorActuatorResponse, state)
		case model.SetDuplicateActuatorRequest:
			state := msg.Payload.(actuator.ActuatorState)
			app.Actuator.SetDuplicateActuatorState(msg.Payload.(actuator.ActuatorState))
			app.sendResponse(model.SetDuplicateActororResponse, state)
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
			// not found in existing DB -> ERROR

			errorQuery := app.ErrorDB.Query(queryString)
			if errorQuery == -1 {
				app.ErrorDB.Insert(queryString, 0)
			}
			go app.Actuator.SetErrorActuatorState(actuator.OnState)
			go app.sendResponse(model.SetErrorActuatorResponse, actuator.OnState)

		} else if existingDBResult == 1 {
			// found barcode
			// do something
			app.ScannedDB.Insert(queryString, 0)
			app.ScannedDB.Query(queryString)
			app.QueryCounter++
			app.TotalCounter++
		} else {
			// found duplicated query
			duplicateQuery := app.DuplicatedItemDB.Query(queryString)

			if duplicateQuery == -1 {
				app.DuplicatedItemDB.Insert(queryString, 0)
			}
			go app.Actuator.SetDuplicateActuatorState(actuator.OnState)
			go app.sendResponse(model.SetDuplicateActororResponse, actuator.OnState)
		}

		if app.QueryCounter == app.QueryCounterLimit {
			app.cleanUp()
		}

		app.sendResponse(model.CurrentCounterUpdateResponse, app.QueryCounter)
		app.sendResponse(model.TotalCounterUpdateResponse, app.TotalCounter)
		log.Printf("Query result %s : %d \n", queryString, existingDBResult)
	}

	defer app.cleanUp()
}
