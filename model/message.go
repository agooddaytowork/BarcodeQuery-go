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
	SetDuplicateActororResponse
	SetErrorActuatorRequest
	SetErrorActuatorResponse
	ResetAppRequest
	RestAppResponse
	SetCurrentCounterLimitRequest
	SetCurrentCounterLimitResponse
	GetNumberOfItemInListRequest
	GetNumberOfItemInListResponse
)

type BarcodeQueryMessage struct {
	MessageType MessageType `json:"message_type"`
	Payload     any         `json:"payload"`
}
