package app

import (
	"BarcodeQuery/reader"
)

type BarcodeAppConfig struct {
	ExistingDBPath                    string            `json:"existing_db_path"`
	ErrorDBPath                       string            `json:"error_db_path"`
	ScannedDBPath                     string            `json:"scanned_db_path"`
	DuplicatedDBPath                  string            `json:"duplicated_db_path"`
	BarcodeReaderType                 reader.ReaderType `json:"barcode_reader_type"`
	BarcodeReaderURI                  string            `json:"barcode_reader_uri"`
	ValidateLotReaderType             reader.ReaderType `json:"validate_lot_reader_type"`
	ValidateLotReaderURI              string            `json:"validate_lot_reader_uri"`
	ReaderDuplicateDebounceIntervalMs int               `json:"reader_duplicate_debounce_interval_ms"`
	QueryCounterLimit                 int               `json:"query_counter_limit"`
	EnableActuator                    bool              `json:"enable_actuator"`
	EnableSoundAlert                  bool              `json:"enable_sound_alert"`
	AutoResetActuator                 bool              `json:"auto_reset_actuator"`
	WebStaticFilePath                 string            `json:"web_static_file_path"`
	ActuatorType                      string            `json:"actuator_type"`
	ActuatorURI                       string            `json:"actuator_uri"`
	DebugMode                         bool              `json:"debug_mode"`
	ValidateLot                       bool              `json:"validate_lot"`
}
