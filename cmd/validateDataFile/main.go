package main

import (
	"BarcodeQuery/classifier"
	"BarcodeQuery/hashing"
	"flag"
	"log"
	"os"
	"strings"
)

func main() {

	filePath := flag.String("f", "data.txt", "File path")
	flag.Parse()
	log.Printf("Đọc file %s \n", *filePath)
	data, err := os.ReadFile(*filePath)

	totalBarcodeDuplicate := 0
	totalSerialDuplicate := 0
	totalHashDuplicate := 0

	if err != nil {
		log.Fatal("lỗi đọc file", err)
	}
	elements := strings.Split(string(data), "\r")

	var serialDupFrequency = make(map[string]int)
	var barcodeDupFrequency = make(map[string]int)
	var hashDupFrequency = make(map[string]int)

	barcodeNSerialClassifier := classifier.BarcodeNSerialTupleClassifier{}

	log.Println("Kiểm tra barcode và serial bị trùng ...")
	log.Println("_______________________________________")
	log.Println("IN LOG LỖI NẾU CÓ")
	for _, line := range elements {
		//log.Printf("____dòng %d _____", index)
		barcode, serial := barcodeNSerialClassifier.Classify(line)

		somethingWrong := false

		if v, ok := serialDupFrequency[serial]; ok {
			somethingWrong = true
			serialDupFrequency[serial] = v + 1
			log.Printf("Serial %s bị trùng lần %d  \n", serial, v)
			totalSerialDuplicate++
		} else {
			serialDupFrequency[serial] = 1
		}

		if v, ok := barcodeDupFrequency[barcode]; ok {
			somethingWrong = true
			barcodeDupFrequency[barcode] = v + 1
			log.Printf("Barcode %s bị trùng lần %d  \n", barcode, v)
			totalBarcodeDuplicate++

		} else {
			barcodeDupFrequency[barcode] = 1
		}

		if !somethingWrong {
			//log.Printf("ok \n")
		}
	}

	log.Println("___________")
	log.Println("Băm mã barcode và kiểm tra bị trùng ...")
	log.Println("IN LOG LỖI NẾU CÓ")
	hasher := hashing.BarcodeSHA256HasherImpl{}
	for barcode := range barcodeDupFrequency {
		hash := hasher.Hash(barcode)
		if v, ok := hashDupFrequency[hash]; ok {
			hashDupFrequency[hash] = v + 1
			log.Printf("Barcode %s, Hash %s bị trùng %d lần \n", barcode, hash, v)
			totalHashDuplicate++
		} else {
			barcodeDupFrequency[barcode] = 1
		}
	}
	log.Println("Đã kiểm tra xong!")
	log.Println("______________________")
	log.Println("______________________")
	log.Println("______________________")
	log.Println("____Tổng Kết ______")
	log.Printf("Số lần barcode bị trùng: %d", totalBarcodeDuplicate)
	log.Printf("Số lần serial bị trùng: %d", totalSerialDuplicate)
	log.Printf("Số lần mã băm bị trùng: %d", totalHashDuplicate)

}
