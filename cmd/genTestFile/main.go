package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Create("test/100k.txt")
	check(err)
	for i := 0; i < 100000; i++ {
		f.WriteString(fmt.Sprintf("%0.12d \n", i))
	}

	defer f.Close()
}
