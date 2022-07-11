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
	ExistingDB         db.BarcodeDB
	DuplicatedItemDB   db.BarcodeDB
	ErrorDB            db.BarcodeDB
	ScannedDB          db.BarcodeDB
	PersistedScannedDB db.BarcodeDB
	Reader             reader.BarcodeReader
	CounterReport      model.CounterReport
	Broadcaster        *broadcast.Broadcaster
	ClientListener     *broadcast.Listener
	Actuator           actuator.BarcodeActuator
	Config             BarcodeAppConfig
}

func (app *BarcodeQueryAppImpl) sendResponse(msgType model.MessageType, payload any) {
	app.Broadcaster.Send(
		model.BarcodeQueryMessage{
			MessageType: msgType,
			Payload:     payload,
		})
}

func (app *BarcodeQueryAppImpl) handleAppReset() {
	app.PersistedScannedDB.Clear()
	app.PersistedScannedDB.Dump()

	app.ExistingDB.Clear()
	app.ExistingDB.Load()
	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
	app.CounterReport.TotalCounter = 0
	app.CounterReport.QueryCounter = 0
	app.CounterReport.PackageCounter = 0
	app.CounterReport.NumberOfItemInExistingDB = app.ExistingDB.GetDBLength()
	app.sendResponse(model.RestAppResponse, "ok")
	app.sendResponse(model.CounterReportResponse, app.CounterReport)
}

func (app *BarcodeQueryAppImpl) handleClientRequest() {
	for true {
		request := <-app.ClientListener.Channel()
		msg := request.(model.BarcodeQueryMessage)
		switch msg.MessageType {
		case model.CurrentCounterUpdateRequest:
			app.sendResponse(model.CurrentCounterUpdateResponse, app.CounterReport.QueryCounter)
		case model.TotalCounterUpdateRequest:
			app.sendResponse(model.TotalCounterUpdateResponse, app.CounterReport.TotalCounter)
		case model.SetErrorActuatorRequest:
			state := actuator.GetState(msg.Payload.(bool))
			app.Actuator.SetErrorActuatorState(state)
			app.sendResponse(model.SetErrorActuatorResponse, state)
		case model.SetDuplicateActuatorRequest:
			state := actuator.GetState(msg.Payload.(bool))
			app.Actuator.SetDuplicateActuatorState(state)
			app.sendResponse(model.SetDuplicateActuatorResponse, state)
		case model.SetCurrentCounterLimitRequest:
			app.CounterReport.QueryCounterLimit = msg.Payload.(int)
			app.sendResponse(model.SetCurrentCounterLimitResponse, msg.Payload.(int))
		case model.ResetAppRequest:
			app.handleAppReset()
		case model.GetNumberOfItemInListRequest:
			app.sendResponse(model.GetNumberOfItemInListResponse, app.CounterReport.NumberOfItemInExistingDB)
		case model.CounterReportRequest:
			app.sendResponse(model.CounterReportResponse, app.CounterReport)
		case model.GetConfigRequest:
			app.sendResponse(model.GetConfigResponse, app.Config)
		case model.SetConfigRequest:
			app.Config = msg.Payload.(BarcodeAppConfig)
			app.sendResponse(model.SetConfigResponse, 1)
		}
	}
}

func (app *BarcodeQueryAppImpl) cleanUp() {
	log.Println("Cleaning up")
	app.sendResponse(model.ResetAllCountersResponse, 0)
	app.CounterReport.QueryCounter = 0
	app.PersistedScannedDB.Dump()
	app.ScannedDB.DumpWithTimeStamp()
	app.ErrorDB.DumpWithTimeStamp()
	app.DuplicatedItemDB.DumpWithTimeStamp()
	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
}

func (app *BarcodeQueryAppImpl) syncPersistedScannedDBToExistingDB() {
	app.ExistingDB.Sync(app.PersistedScannedDB.GetStore())
	app.CounterReport.TotalCounter = app.PersistedScannedDB.GetDBLength()
}

func (app *BarcodeQueryAppImpl) Run() {
	app.CounterReport.NumberOfItemInExistingDB = app.ExistingDB.GetDBLength()
	app.syncPersistedScannedDBToExistingDB()
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
			app.PersistedScannedDB.Insert(queryString, 0)
			app.ScannedDB.Query(queryString)
			app.CounterReport.QueryCounter++
			app.CounterReport.TotalCounter++
		} else {
			// found duplicated query
			duplicateQuery := app.DuplicatedItemDB.Query(queryString)
			if duplicateQuery == -1 {
				app.DuplicatedItemDB.Insert(queryString, 0)
			}
			go app.Actuator.SetDuplicateActuatorState(actuator.OnState)
			go app.sendResponse(model.SetDuplicateActuatorResponse, actuator.OnState)
		}
		if app.CounterReport.QueryCounter == app.CounterReport.QueryCounterLimit {
			app.CounterReport.PackageCounter++
			app.cleanUp()
		}
		app.sendResponse(model.CounterReportResponse, app.CounterReport)
		log.Printf("Query result %s : %d \n", queryString, existingDBResult)
	}

	defer app.cleanUp()
}
