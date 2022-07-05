package reader

import (
	"bufio"
	"log"
	"net"
	"time"
)

type TCPReader struct {
	URL           string
	client        *net.Conn
	SpawnedThread bool
	ReportChannel chan string
}

func (r *TCPReader) connect() {

	try := true
	for try {
		if r.client == nil {
			log.Printf("Creating connection to %s \n", r.URL)
			conn, err := net.Dial("tcp", r.URL)
			if err != nil {
				log.Println(err)
				log.Println("Sleep 5 second then retry to connect to host")
				time.Sleep(time.Second * 5)

			} else {
				r.client = &conn
			}
		} else {
			try = false
		}
	}
}

func (r *TCPReader) readAsync() {
	r.connect()

	run := true

	for run {
		log.Printf("Wait for input from %s ...", r.URL)
		status, err := bufio.NewReader(*r.client).ReadString('\n')
		if err != nil {
			log.Println(err)
			r.client = nil
			r.SpawnedThread = false
			run = false
			r.ReportChannel <- ""
		} else {
			r.ReportChannel <- status
		}
	}
}

func (r *TCPReader) Read() string {

	if r.SpawnedThread == false {
		r.SpawnedThread = true
		go r.readAsync()
	}

	status := <-r.ReportChannel
	return status
}
