package model

type MessageType int16

const (
	DBQueryNoti MessageType = iota
	CounterNoti
)

type BarcodeQueryMessage struct {
	MessageType MessageType `json:"message_type"`
	Payload     any         `json:"payload"`
}
