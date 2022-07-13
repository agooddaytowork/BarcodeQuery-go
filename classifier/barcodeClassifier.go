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
	return strings.Trim(multiStrings[1], " \r\n\t")
}
