package reader

import (
	"os"
	"strings"
	"time"
)

type TestFileReader struct {
	Interval time.Duration
	Store    []string
	Index    int
}

func (r *TestFileReader) Load(filePath string) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	r.Store = strings.Split(string(data), "\n")
	r.Index = -1

}

func (r *TestFileReader) Read() string {
	time.Sleep(r.Interval)
	r.Index++

	if r.Index > len(r.Store)-1 {
		r.Index = 0
	}
	return strings.Trim(r.Store[r.Index], " ")
}
