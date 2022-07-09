package main

import (
	"BarcodeQuery/app"
	"BarcodeQuery/web"
	"flag"
	"github.com/textileio/go-threads/broadcast"
)

func main() {

	configPath := flag.String("c", "C:\\Users\\tam\\Workspace\\BarcodeQuery-go\\target\\config.json", "Config path")
	flag.Parse()

	theConfig := app.LoadConfigFromFile(*configPath)
	dbBroadCast := broadcast.NewBroadcaster(100)
	clientBroadCast := broadcast.NewBroadcaster(100)

	program := app.GetBarcodeQueryAppImpl(theConfig, dbBroadCast, clientBroadCast)
	theWeb := web.GetBarcodeQueryWebImpl(dbBroadCast, clientBroadCast, theConfig.WebStaticFilePath)

	go theWeb.Run()
	program.Run()

}
