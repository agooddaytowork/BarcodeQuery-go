package web

import "BarcodeQuery/app"

type BarcodeWebApp interface {
	Run()
	RegisterDBCallBack(dbRole app.DBRole, callback func())
}
