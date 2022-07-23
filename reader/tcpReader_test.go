package reader

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTcpReader(t *testing.T) {

	listen, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}

	go func() {
		conn, err := listen.Accept()
		if err != nil {
			return
		}

		for i := 0; i < 100; i++ {
			conn.Write([]byte(fmt.Sprintf("%0.12d\n", i)))
			time.Sleep(50 * time.Millisecond)
		}
		conn.Close()
	}()

	reportChannel := make(chan string, 1000)

	tcpReader := TCPReader{
		URL:           "localhost:9000",
		client:        nil,
		SpawnedThread: false,
		ReportChannel: reportChannel,
	}

	result := make(map[string]int)

	i := 0
	for i < 100 {
		data := tcpReader.Read()

		if value, found := result[data]; found {
			fmt.Printf("Found duplicate %s", data)
			result[data] = value + 1
		} else {
			result[data] = 1
		}
		i++
		fmt.Println(i)
	}
}
