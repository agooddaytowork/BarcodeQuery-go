package classifier

type BarcodeTupleClassifier struct {
}

func (c *BarcodeTupleClassifier) Classify(input string) (string, string) {
	barcodeClassifier := barcodeClassifier{}
	return barcodeClassifier.Classify(input), ""
}
