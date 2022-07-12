package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func LoadConfigFromFile(filePath string, output any) {
	file, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal([]byte(file), &output)
	if err != nil {
		log.Panicf("File config %s không đúng format, lỗi: %s \n", filePath, err.Error())
	}
}

func DumpConfigToFile(filePath string, config any) {
	file, _ := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	jsonString, _ := json.MarshalIndent(config, "", "    ")
	file.Write(jsonString)
	defer file.Close()
}
