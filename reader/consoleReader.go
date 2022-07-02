package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ConsoleReader struct {
}

func (r *ConsoleReader) Read() string {
	fmt.Println("Please input a query: ")
	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	//fmt.Println("you have entered ", text)

	return strings.Trim(text, " \n\t")
}
