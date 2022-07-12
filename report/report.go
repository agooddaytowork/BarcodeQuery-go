package report

type Report interface {
	WriteRecord(queryString string, value int)
	Flush()
}
