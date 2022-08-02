package app

import (
	"BarcodeQuery/actuator"
	"BarcodeQuery/model"
	"BarcodeQuery/util"
	"encoding/json"
)

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
