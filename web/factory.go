package web

import "github.com/textileio/go-threads/broadcast"

func GetBarcodeQueryWebImpl(dbBroadCast *broadcast.Broadcaster, clientBroadCast *broadcast.Broadcaster, webStaticPath string, barcodeListFilePath string) BarcodeQueryWebImpl {
	return BarcodeQueryWebImpl{
		Broadcaster:         dbBroadCast,
		ClientBroadCast:     clientBroadCast,
		StaticFilePath:      webStaticPath,
		BarcodeListFilePath: barcodeListFilePath,
	}
}
