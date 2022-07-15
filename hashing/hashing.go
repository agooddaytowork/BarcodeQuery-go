package hashing

type BarcodeHashser interface {
	Hash(input string) string
	CheckHash(input, hash string) bool
}
