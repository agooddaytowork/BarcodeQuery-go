package model

type DebugRecord struct {
	Serial           string `json:"serial"`
	Hash             string `json:"hash"`
	Barcode          string `json:"barcode"`
	ScannedTimestamp int64  `json:"scanned_timestamp"`
	State            string `json:"state"`
}
