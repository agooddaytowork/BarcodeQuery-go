package model

type PersistedSerialRecord struct {
	Serial           string `json:"serial"`
	ScannedTimestamp int64  `json:"scanned_timestamp"`
	Lot              string `json:"lot"`
}
