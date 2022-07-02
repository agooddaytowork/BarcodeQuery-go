package main

import (
	"fmt"
	"time"
)

func main() {

	for i := 0; i < 100; i++ {
		fmt.Printf("%0.12d \n", i)
		time.Sleep(time.Millisecond * 500)
	}

}
