package reader

import (
	"fmt"
	"time"
)

func NewReader(readerType ReaderType, readerURI string, debounceInterval int) BarcodeReader {
	var barcodeReader BarcodeReader

	switch readerType {
	case TestFileReaderType:
		testFileReader := TestFileReader{
			Interval: time.Millisecond * 200,
		}
		testFileReader.Load(readerURI)
		barcodeReader = &testFileReader

	case ConsoleReaderType:
		barcodeReader = &ConsoleReader{}

	case TCPReaderType:
		barcodeReader = &TCPReader{
			URL:                         readerURI,
			SpawnedThread:               false,
			ReportChannel:               make(chan string, 1000),
			DuplicateDebounceIntervalMs: debounceInterval,
		}
	default:
		panic(fmt.Sprintf("Unsupported reader, only support %s/%s/%s", TestFileReaderType, ConsoleReaderType, TCPReaderType))
	}

	return barcodeReader
}
