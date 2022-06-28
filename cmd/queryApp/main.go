package main

import (
	app2 "BarcodeQuery/app"
	"BarcodeQuery/db"
	"BarcodeQuery/reader"
)

func main() {
	theDB := db.BarcodeDBHashStorageImpl{
		FilePath: "/Users/tam/Workspace/Duc/BarcodeQuery/test/100k.txt",
		Store:    make(map[string]int),
	}
	err := theDB.Load()

	if err != nil {
		panic(err)
	}
	program := app2.BarcodeQueryAppImpl{
		Db:     &theDB,
		Reader: &reader.ConsoleReader{},
	}

	program.Run()

}
