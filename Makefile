
all: query-app file-validator
query-app:
	GOOS=windows GOARCH=amd64 go build -o target/barcodequery.exe ./cmd/queryApp/main.go
file-validator:
	GOOS=windows GOARCH=amd64 go build -o target/fileValidator.exe ./cmd/validateDataFile/main.go