package classifier

type SerialNBarcodeTupleClassifier struct {
}

func (c *SerialNBarcodeTupleClassifier) Classify(input string) (string, string) {
	serialClassifier := serialClassifier{}
	barcodeClassifier := barcodeClassifier{}
	return serialClassifier.Classify(input), barcodeClassifier.Classify(input)
}
