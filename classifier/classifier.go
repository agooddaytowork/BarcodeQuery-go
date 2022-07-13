package classifier

type classifier interface {
	Classify(input string) string
}

type TupleClassifier interface {
	Classify(input string) (string, string)
}
