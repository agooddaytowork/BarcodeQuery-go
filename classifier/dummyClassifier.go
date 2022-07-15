package classifier

import (
	"strings"
)

type dummyClassifier struct {
}

func (c *dummyClassifier) Classify(input string) string {
	return strings.Trim(input, " \r\n\t")
}
