package reader

import (
	"log"
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
	r.Index = 0

}

func (r *TestFileReader) Read() string {
	time.Sleep(r.Interval)
	r.Index++

	if r.Index == len(r.Store) {
		log.Println("Đã đọc hết test file, thoát chương trình")
		os.Exit(0)
	}
	return strings.Trim(r.Store[r.Index], " \r")
}
