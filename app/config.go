package app

import (
	"BarcodeQuery/reader"
)

type BarcodeAppConfig struct {
	ExistingDBPath    string            `json:"existing_db_path"`
	ErrorDBPath       string            `json:"error_db_path"`
	ScannedDBPath     string            `json:"scanned_db_path"`
	DuplicatedDBPath  string            `json:"duplicated_db_path"`
	ReaderType        reader.ReaderType `json:"reader_type"`
	ReaderURI         string            `json:"reader_uri"`
	QueryCounterLimit int               `json:"query_counter_limit"`
	EnableActuator    bool              `json:"enable_actuator"`
	EnableSoundAlert  bool              `json:"enable_sound_alert"`
	AutoResetActuator bool              `json:"auto_reset_actuator"`
	WebStaticFilePath string            `json:"web_static_file_path"`
}
