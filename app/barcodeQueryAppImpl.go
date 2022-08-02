package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/classifier"
	"BarcodeQuery/db"
	"BarcodeQuery/hashing"
	"BarcodeQuery/model"
	"BarcodeQuery/reader"
	"BarcodeQuery/util"
	"errors"
	"fmt"
	"github.com/textileio/go-threads/broadcast"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BarcodeQueryAppImpl struct {
	BarcodeExistingDB           db.SerialDB
	SerialAndBarcodeDB          db.SerialNBarcodeDB
	BarcodeAndSerialDB          db.SerialNBarcodeDB
	DuplicatedItemDB            db.SerialDB
	DebugDB                     db.SerialDB
	ErrorDB                     db.SerialDB
	ScannedDB                   db.SerialDB
	PersistedScannedDB          db.PersistedSerialRecordDB
	MainBarcodeReader           reader.BarcodeReader
	ValidateLotBarcodeReader    reader.BarcodeReader
	CounterReport               model.CounterReport
	Broadcaster                 *broadcast.Broadcaster
	ClientListener              *broadcast.Listener
	Actuator                    actuator.BarcodeActuator
	Config                      BarcodeAppConfig
	ConfigPath                  string
	Hasher                      hashing.BarcodeHashser
	TestMode                    bool
	RunMode                     string
	ValidateModeTerminateSignal chan interface{}
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

func (app *BarcodeQueryAppImpl) mainLogic() {
	run := true
	var debugArray []model.DebugRecord

	for run {
		barcode := app.MainBarcodeReader.Read()
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
			var persistedRecord model.PersistedSerialRecord
			if record, ok := app.PersistedScannedDB.Query(serialNumber); ok {
				persistedRecord = record
			} else {
				persistedRecord = model.PersistedSerialRecord{
					Serial:           serialNumber,
					ScannedTimestamp: time.Now().Unix(),
					Lot:              "Lô hiện tại",
				}
			}

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
		app.sendResponse(model.CounterReportResponse, app.CounterReport)

		if app.CounterReport.QueryCounter == app.CounterReport.QueryCounterLimit {
			if debugArray != nil && len(debugArray) > 0 {
				util.DumpConfigToFile("debug/debug-"+strconv.FormatInt(time.Now().Unix(), 10)+".json", debugArray)
				debugArray = nil
			}
			app.sendResponse(model.CurrentCounterHitLimitNoti, model.CounterHitLimitPayload{
				LotIdentifier: app.getLotIdentifier(),
			})
			if app.Config.ValidateLot {
				app.runValidateMode()
			}
		}
		log.Printf("Query result %s : %d \n", barcodeHash, existingDBResult)
	}
	util.DumpConfigToFile("debug/debug-"+strconv.FormatInt(time.Now().Unix(), 10)+".json", debugArray)

}

//TODO add an option for retrying validate mode
func (app *BarcodeQueryAppImpl) runValidateMode() {

	app.sendResponse(model.ValidateLotStartNoti, 1)
	//go func() {
	//	<-app.ValidateModeTerminateSignal
	//	running = false
	//}()
	validateString := app.ValidateLotBarcodeReader.Read()
	app.sendResponse(model.ValidateStringResponse, validateString)
	unexpectedSerialNumbers, _ := app.validateLot(validateString)

	app.sendResponse(model.ValidateLotResponse, model.ValidateLotResponsePayload{
		UnexpectedSerialNumbers: unexpectedSerialNumbers,
	})

}

func (app *BarcodeQueryAppImpl) validateLot(validateString string) ([]string, error) {

	elements := strings.Split(validateString, "-")

	from, err := strconv.Atoi(elements[0])

	if err != nil {
		return nil, err
	}

	to, err := strconv.Atoi(elements[1])
	if err != nil {
		return nil, err
	}
	scannedResult := app.ScannedDB.GetStore()
	var unexpectedSerialNumbers []string
	for i := from; i <= to; i++ {
		serialNumber := fmt.Sprintf("%0.12d", i)
		if _, found := scannedResult[serialNumber]; !found {
			unexpectedSerialNumbers = append(unexpectedSerialNumbers, serialNumber)
		}
	}

	if len(unexpectedSerialNumbers) > 0 {
		sort.Strings(unexpectedSerialNumbers)

		return unexpectedSerialNumbers, errors.New("found unexpected serial number")
	}

	return unexpectedSerialNumbers, nil
}

func (app *BarcodeQueryAppImpl) Run() {

	app.CounterReport.NumberOfItemInExistingDB = app.BarcodeExistingDB.GetDBLength()
	app.syncPersistedScannedDBToExistingDB()

	go app.handleClientRequest()
	go app.DuplicatedItemDB.HandleClientRequest()
	go app.ErrorDB.HandleClientRequest()
	go app.ScannedDB.HandleClientRequest()
	app.mainLogic()
	defer func() {
		app.cleanUp()
	}()
}
