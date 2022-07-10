package reader

type BarcodeReader interface {
	Read() string
}

type ReaderType string

const (
	TestFileReaderType ReaderType = "test_file"
	ConsoleReaderType  ReaderType = "console"
	TCPReaderType      ReaderType = "tcp_reader"
)
