package classifier

import (
	"fmt"
	"strings"
)

type serialClassifier struct {
}

func (c *serialClassifier) Classify(input string) string {

	if input == "" {
		return input
	}
	multiStrings := strings.Split(input, "\t")

	if len(multiStrings) < 2 {
		panic(fmt.Sprintf("invalid input %s, SerialClassifier input must be <serialnumber>\t<barcode>", input))
	}
	return strings.Trim(multiStrings[0], " \r\n\t")
}
