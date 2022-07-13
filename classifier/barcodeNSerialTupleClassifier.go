package classifier

type BarcodeNSerialTupleClassifier struct {
}

func (c *BarcodeNSerialTupleClassifier) Classify(input string) (string, string) {
	serialClassifier := serialClassifier{}
	barcodeClassifier := barcodeClassifier{}
	return barcodeClassifier.Classify(input), serialClassifier.Classify(input)
}
