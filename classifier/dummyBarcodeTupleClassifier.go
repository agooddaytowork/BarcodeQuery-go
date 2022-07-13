package classifier

type DummyBarcodeTupleClassifier struct {
}

func (c *DummyBarcodeTupleClassifier) Classify(input string) (string, string) {
	dummyClassifier := dummyClassifier{}
	return dummyClassifier.Classify(input), ""
}
