package reader

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type ConsoleReader struct {
}

func (r *ConsoleReader) Read() string {
	log.Println("Please input a query: ")
	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	//fmt.Println("you have entered ", text)

	return strings.Trim(text, " \n\t")
}
