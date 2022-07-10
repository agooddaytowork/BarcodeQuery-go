package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/db"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"github.com/textileio/go-threads/broadcast"
	"log"
)

type BarcodeQueryAppImpl struct {
	ExistingDB               db.BarcodeDB
	DuplicatedItemDB         db.BarcodeDB
	ErrorDB                  db.BarcodeDB
	ScannedDB                db.BarcodeDB
	Reader                   reader.BarcodeReader
	QueryCounter             int
	QueryCounterLimit        int
	TotalCounter             int
	NumberOfItemInExistingDB int
	Broadcaster              *broadcast.Broadcaster
	ClientListener           *broadcast.Listener
	Actuator                 actuator.BarcodeActuator
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
			state := actuator.GetState(msg.Payload.(bool))
			app.Actuator.SetErrorActuatorState(state)
			app.sendResponse(model.SetErrorActuatorResponse, state)
		case model.SetDuplicateActuatorRequest:
			state := actuator.GetState(msg.Payload.(bool))
			app.Actuator.SetDuplicateActuatorState(state)
			app.sendResponse(model.SetDuplicateActuatorResponse, state)
		case model.SetCurrentCounterLimitRequest:
			app.QueryCounterLimit = msg.Payload.(int)
			app.sendResponse(model.SetCurrentCounterLimitResponse, msg.Payload.(int))
		case model.ResetAppRequest:
			// todo: handle reset request
			app.sendResponse(model.RestAppResponse, "ok")
		case model.GetNumberOfItemInListRequest:
			app.sendResponse(model.GetNumberOfItemInListResponse, app.NumberOfItemInExistingDB)
		}

	}
}

func (app *BarcodeQueryAppImpl) cleanUp() {
	log.Println("Cleaning up")
	app.sendResponse(model.ResetAllCountersResponse, 0)
	app.QueryCounter = 0
	app.ScannedDB.DumpWithTimeStamp()
	app.ErrorDB.DumpWithTimeStamp()
	app.DuplicatedItemDB.DumpWithTimeStamp()
	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
}

func (app *BarcodeQueryAppImpl) Run() {
	app.NumberOfItemInExistingDB = app.ExistingDB.GetDBLength()
	run := true

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
			go app.sendResponse(model.SetDuplicateActuatorResponse, actuator.OnState)
		}

		app.sendResponse(model.CurrentCounterUpdateResponse, app.QueryCounter)
		app.sendResponse(model.TotalCounterUpdateResponse, app.TotalCounter)
		if app.QueryCounter == app.QueryCounterLimit {
			app.cleanUp()
		}

		log.Printf("Query result %s : %d \n", queryString, existingDBResult)
	}

	defer app.cleanUp()
}
