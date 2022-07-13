package classifier

import (
	"fmt"
	"strings"
)

type barcodeClassifier struct {
}

func (c *barcodeClassifier) Classify(input string) string {
	if input == "" {
		return input
	}
	multiStrings := strings.Split(input, "\t")
	if len(multiStrings) != 2 {
		panic(fmt.Sprintf("invalid input %s, SerialClassifier input must be <serialnumber>\t<barcode>"))
	}
	output := strings.ReplaceAll(multiStrings[1], "\r", "")
	output = strings.ReplaceAll(output, " ", "")
	return strings.Trim(output, " \r\t\n")
}
