package app

import (
	"encoding/json"
	"io/ioutil"
)

type BarcodeAppConfig struct {
	ExistingDBPath    string `json:"existing_db_path"`
	ErrorDBPath       string `json:"error_db_path"`
	ScannedDBPath     string `json:"scanned_db_path"`
	DuplicatedDBPath  string `json:"duplicated_db_path"`
	ReaderType        string `json:"reader_type"`
	ReaderURI         string `json:"reader_uri"`
	QueryCounterLimit int    `json:"query_counter_limit"`
	EnableActuator    bool   `json:"enable_actuator"`
	AutoResetActuator bool   `json:"auto_reset_actuator"`
}

func LoadConfigFromFile(filePath string) BarcodeAppConfig {
	file, _ := ioutil.ReadFile(filePath)
	data := BarcodeAppConfig{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}
