package model

type MessageType int16

const (
	DBQueryNoti MessageType = iota
	DBStateUpdateResponse
	DBStateUpdateRequest
	CurrentCounterUpdateRequest
	CurrentCounterUpdateResponse
	TotalCounterUpdateRequest
	TotalCounterUpdateResponse
	SetDuplicateActuatorRequest
	SetDuplicateActuatorResponse
	SetErrorActuatorRequest
	SetErrorActuatorResponse
	ResetAppRequest
	RestAppResponse
	SetCurrentCounterLimitRequest
	SetCurrentCounterLimitResponse
	GetNumberOfItemInListRequest
	GetNumberOfItemInListResponse
	ResetAllCountersRequest
	ResetAllCountersResponse
	CounterReportRequest
	CounterReportResponse
	GetConfigRequest
	GetConfigResponse
	SetConfigRequest
	SetConfigResponse
)

type BarcodeQueryMessage struct {
	MessageType MessageType `json:"message_type"`
	Payload     any         `json:"payload"`
}
