package app

import (
	"BarcodeQuery/reader"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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

func LoadConfigFromFile(filePath string) BarcodeAppConfig {
	file, _ := ioutil.ReadFile(filePath)
	data := BarcodeAppConfig{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Panicf("File config %s không đúng format, lỗi: %s \n", filePath, err.Error())
	}
	return data
}

func DumpConfigToFile(filePath string, config BarcodeAppConfig) {
	file, _ := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	jsonString, _ := json.MarshalIndent(config, "", "    ")
	file.Write(jsonString)
	defer file.Close()
}
