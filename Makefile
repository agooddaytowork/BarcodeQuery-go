
build-windows:
	GOOS=windows GOARCH=amd64 go build -o target/barcodequery.exe ./cmd/queryApp/main.go
