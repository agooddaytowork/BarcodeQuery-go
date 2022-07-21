package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/classifier"
	"BarcodeQuery/db"
	"BarcodeQuery/hashing"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"BarcodeQuery/util"
	"encoding/json"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"log"
	"sort"
	"strconv"
	"time"
)

type BarcodeQueryAppImpl struct {
	BarcodeExistingDB  db.SerialDB
	SerialAndBarcodeDB db.SerialNBarcodeDB
	BarcodeAndSerialDB db.SerialNBarcodeDB
	DuplicatedItemDB   db.SerialDB
	DebugDB            db.SerialDB
	ErrorDB            db.SerialDB
	ScannedDB          db.SerialDB
	PersistedScannedDB db.PersistedSerialRecordDB
	Reader             reader.BarcodeReader
	CounterReport      model.CounterReport
	Broadcaster        *broadcast.Broadcaster
	ClientListener     *broadcast.Listener
	Actuator           actuator.BarcodeActuator
	Config             BarcodeAppConfig
	ConfigPath         string
	Hasher             hashing.BarcodeHashser
	TestMode           bool
}

func (app *BarcodeQueryAppImpl) sendResponse(msgType model.MessageType, payload any) {
	app.Broadcaster.Send(
		model.BarcodeQueryMessage{
			MessageType: msgType,
			Payload:     payload,
		})
}

func (app *BarcodeQueryAppImpl) handleAppReset() {
	app.BarcodeExistingDB.Clear()
	app.BarcodeExistingDB.Load(&classifier.BarcodeTupleClassifier{})
	app.SerialAndBarcodeDB.Clear()
	app.SerialAndBarcodeDB.Load(&classifier.SerialNBarcodeTupleClassifier{})
	app.BarcodeAndSerialDB.Clear()
	app.BarcodeAndSerialDB.Load(&classifier.BarcodeNSerialTupleClassifier{})
	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
	app.CounterReport.TotalCounter = 0
	app.CounterReport.QueryCounter = 0
	app.CounterReport.PackageCounter = 0
	app.CounterReport.NumberOfCameraScanError = 0
	app.CounterReport.NumberOfItemInExistingDB = app.BarcodeExistingDB.GetDBLength()
	app.sendResponse(model.ResetAppResponse, "ok")
	app.sendResponse(model.CounterReportResponse, app.CounterReport)
	app.syncPersistedScannedDBToExistingDB()
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
			var newConfig BarcodeAppConfig
			jsonString, _ := json.Marshal(msg.Payload)
			json.Unmarshal(jsonString, &newConfig)
			app.Config = newConfig
			app.CounterReport.QueryCounterLimit = app.Config.QueryCounterLimit
			app.sendResponse(model.GetConfigResponse, app.Config)
			app.sendResponse(model.CounterReportResponse, app.CounterReport)
			util.DumpConfigToFile(app.ConfigPath, app.Config)
		case model.CloseCurrentLotRequest:
			app.CounterReport.QueryCounter = 0
			app.CounterReport.PackageCounter++
			app.syncScannedDataToPersistedStorage()
			app.cleanUp()
			app.sendResponse(model.CounterReportResponse, app.CounterReport)
		case model.ResetCurrentCounterRequest:
			app.CounterReport.QueryCounter = 0
			//app.syncScannedDataToPersistedStorage()
			app.handleCurrentCounterReset()
			app.sendResponse(model.CounterReportResponse, app.CounterReport)
		case model.SetCameraErrorActuatorRequest:
			state := actuator.GetState(msg.Payload.(bool))
			app.Actuator.SetCameraErrorActuatorState(state)
			app.sendResponse(model.SetCameraErrorActuatorResponse, state)
		// todo , add camera error actuator
		case model.ResetPersistedFileRequest:
			if !app.TestMode {
				app.PersistedScannedDB.Clear()
				app.PersistedScannedDB.Dump()
			}
			app.handleAppReset()
			app.sendResponse(model.ResetPersistedFileResponse, 1)

		case model.GetDuplicatedItemsStateRequest:
			var duplicatedItemsExistInPersistedRecord []model.PersistedSerialRecord
			for v := range app.DuplicatedItemDB.GetStore() {
				if persistedRecord, ok := app.PersistedScannedDB.Query(v); ok {
					duplicatedItemsExistInPersistedRecord = append(duplicatedItemsExistInPersistedRecord, persistedRecord)
				}
			}
			app.sendResponse(model.GetDuplicatedItemsStateResponse, duplicatedItemsExistInPersistedRecord)

		case model.SetTestModeRequest:
			app.TestMode = msg.Payload.(bool)
			app.sendResponse(model.SetTestModeResponse, app.TestMode)

		case model.GetTestModeStatusRequest:
			app.sendResponse(model.GetTestModeStatusResponse, app.TestMode)

		}
	}
}

// This function is only valid when counter limit is hit
func (app *BarcodeQueryAppImpl) getLotIdentifier() string {
	data := app.ScannedDB.GetStoreAsQueryResultArray()
	sort.Slice(data, func(i, j int) bool {
		return data[i].QueryString < data[j].QueryString
	})

	if len(data) == 0 {
		return ""
	}

	start, _ := strconv.Atoi(data[0].QueryString)
	stop, _ := strconv.Atoi(data[len(data)-1].QueryString)
	return fmt.Sprintf("%d-%d", start, stop)
}

func (app *BarcodeQueryAppImpl) syncScannedDataToPersistedStorage() {
	// get lot number
	lotIdentifier := app.getLotIdentifier()
	log.Printf("lotIdentifier: %s", lotIdentifier)

	for serialNumber := range app.ScannedDB.GetStore() {
		app.PersistedScannedDB.Insert(serialNumber, model.PersistedSerialRecord{
			Serial:           serialNumber,
			ScannedTimestamp: time.Now().Unix(),
			Lot:              lotIdentifier,
		})
	}
}

func (app *BarcodeQueryAppImpl) cleanUp() {
	log.Println("Cleaning up")
	app.sendResponse(model.ResetAllCountersResponse, 0)
	app.CounterReport.QueryCounter = 0

	if !app.TestMode {
		app.PersistedScannedDB.Dump()
		app.ScannedDB.DumpWithTimeStamp()
		app.ErrorDB.DumpWithTimeStamp()
		app.DuplicatedItemDB.DumpWithTimeStamp()
	}

	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
}

func (app *BarcodeQueryAppImpl) handleCurrentCounterReset() {
	app.CounterReport.QueryCounter = 0
	app.BarcodeExistingDB.Clear()
	app.BarcodeExistingDB.Load(&classifier.BarcodeTupleClassifier{})
	app.SerialAndBarcodeDB.Clear()
	app.SerialAndBarcodeDB.Load(&classifier.SerialNBarcodeTupleClassifier{})
	app.BarcodeAndSerialDB.Clear()
	app.BarcodeAndSerialDB.Load(&classifier.BarcodeNSerialTupleClassifier{})
	app.ScannedDB.Clear()
	app.ErrorDB.Clear()
	app.DuplicatedItemDB.Clear()
	app.syncPersistedScannedDBToExistingDB()
	app.sendResponse(model.ResetCurrentCounterResponse, 0)
}

func (app *BarcodeQueryAppImpl) syncPersistedScannedDBToExistingDB() {
	theMap := make(map[string]int)
	serialNBarcodeMap := app.SerialAndBarcodeDB.GetStore()
	for serial := range app.PersistedScannedDB.GetStore() {
		barcode := serialNBarcodeMap[serial]
		theMap[barcode] = 1
	}
	app.BarcodeExistingDB.Sync(theMap)
	app.CounterReport.TotalCounter = app.PersistedScannedDB.GetDBLength()
}

func (app *BarcodeQueryAppImpl) Run() {

	var debugArray []model.DebugRecord

	app.CounterReport.NumberOfItemInExistingDB = app.BarcodeExistingDB.GetDBLength()
	app.syncPersistedScannedDBToExistingDB()
	run := true

	go app.handleClientRequest()
	go app.DuplicatedItemDB.HandleClientRequest()
	go app.ErrorDB.HandleClientRequest()
	go app.ScannedDB.HandleClientRequest()

	for run {
		barcode := app.Reader.Read()
		barcodeHash := app.Hasher.Hash(barcode)

		if barcode == "" {
			continue
		}

		if barcode == CAMERA_ERROR_1 {

			if app.Config.DebugMode {
				debugArray = append(debugArray, model.DebugRecord{
					Serial:           CAMERA_ERROR_1,
					Hash:             "N/A",
					ScannedTimestamp: time.Now().Unix(),
					Barcode:          "",
					State:            "camera_error",
				})
			}

			app.Actuator.SetCameraErrorActuatorState(actuator.OnState)
			app.CounterReport.NumberOfCameraScanError++
			app.sendResponse(model.SetCameraErrorActuatorResponse, true)
			app.sendResponse(model.CounterReportResponse, app.CounterReport)
			continue
		}
		existingDBResult := app.BarcodeExistingDB.Query(barcodeHash)
		if existingDBResult < 0 {
			// not found in existing DB -> ERROR
			errorQuery := app.ErrorDB.Query(barcode)
			if errorQuery == -1 {
				app.ErrorDB.Insert(barcode, 0)
			}
			go app.Actuator.SetErrorActuatorState(actuator.OnState)
			go app.sendResponse(model.SetErrorActuatorResponse, actuator.OnState)

			if app.Config.DebugMode {
				debugArray = append(debugArray, model.DebugRecord{
					Serial:           "N/A",
					Barcode:          barcode,
					Hash:             "N/A",
					ScannedTimestamp: time.Now().Unix(),
					State:            "not found in existing db",
				})
			}
		} else if existingDBResult == 1 {
			// found barcode
			// do something

			serialNumber := app.BarcodeAndSerialDB.Query(barcodeHash)

			app.ScannedDB.Insert(serialNumber, 0)
			//app.PersistedScannedDB.Insert(serialNumber, 0)
			app.ScannedDB.Query(serialNumber)
			app.CounterReport.QueryCounter++
			app.CounterReport.TotalCounter++

			if app.Config.DebugMode {
				debugArray = append(debugArray, model.DebugRecord{
					Serial:           serialNumber,
					Barcode:          barcode,
					Hash:             barcodeHash,
					ScannedTimestamp: time.Now().Unix(),
					State:            "Found in existing DB",
				})
			}
		} else {
			// found duplicated query
			serialNumber := app.BarcodeAndSerialDB.Query(barcodeHash)
			duplicateQuery := app.DuplicatedItemDB.Query(serialNumber)
			persistedRecord, _ := app.PersistedScannedDB.Query(serialNumber)

			if duplicateQuery == -1 {
				app.DuplicatedItemDB.Insert(serialNumber, 0)
			}

			go app.Actuator.SetDuplicateActuatorState(actuator.OnState)
			go app.sendResponse(model.SetDuplicateActuatorResponse, actuator.OnState)
			go app.sendResponse(model.DuplicatedItemNoti, persistedRecord)
			app.CounterReport.QueryCounter++
			app.CounterReport.TotalCounter++

			if app.Config.DebugMode {
				debugArray = append(debugArray, model.DebugRecord{
					Serial:           serialNumber,
					Barcode:          barcode,
					Hash:             barcodeHash,
					ScannedTimestamp: time.Now().Unix(),
					State:            "Duplicated",
				})
			}
		}
		if app.CounterReport.QueryCounter == app.CounterReport.QueryCounterLimit {
			app.sendResponse(model.CurrentCounterHitLimitNoti, model.CounterHitLimitPayload{
				LotIdentifier: app.getLotIdentifier(),
			})
			util.DumpConfigToFile("debug/debug-"+strconv.FormatInt(time.Now().Unix(), 10)+".json", debugArray)
		}
		app.sendResponse(model.CounterReportResponse, app.CounterReport)
		log.Printf("Query result %s : %d \n", barcodeHash, existingDBResult)
	}

	defer func() {
		util.DumpConfigToFile("debug/debug-"+strconv.FormatInt(time.Now().Unix(), 10)+".json", debugArray)
		app.cleanUp()
	}()
}
