package report

import (
	"encoding/csv"
	"fmt"
	"time"
)

type CSVReport struct {
	FilePath  string
	csvWriter *csv.Writer
}

func (report *CSVReport) WriteRecord(queryString string, value int) {
	if report.csvWriter == nil {
		return
	}
	report.csvWriter.Write([]string{
		queryString,
		fmt.Sprintf("%d", value),
		time.Now().Format("2006-01-02-15-04-05"),
	})
}
func (report *CSVReport) Flush() {
	if report.csvWriter == nil {
		return
	}
	report.csvWriter.Flush()
}
