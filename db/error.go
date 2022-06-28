package db

type BarcodeDBError struct {
	ExceptionMsg string
}

func (m *BarcodeDBError) Error() string {
	return m.Error()
}
