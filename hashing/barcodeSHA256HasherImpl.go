package hashing

import (
	"crypto/sha256"
	"fmt"
)

type BarcodeSHA256HasherImpl struct {
}

func (h *BarcodeSHA256HasherImpl) Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

func (h *BarcodeSHA256HasherImpl) CheckHash(input, hash string) bool {
	return h.Hash(input) == hash
}
