package web

import "github.com/textileio/go-threads/broadcast"

func GetBarcodeQueryWebImplementation(dbBroadCast *broadcast.Broadcaster, clientBroadCast *broadcast.Broadcaster, webStaticPath string) BarcodeQueryWebImpl {
	return BarcodeQueryWebImpl{
		Broadcaster:     dbBroadCast,
		ClientBroadCast: clientBroadCast,
		StaticFilePath:  webStaticPath,
	}
}
