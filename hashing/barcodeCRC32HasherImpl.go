package hashing

import "hash/crc32"

type BarcodeCRC32HasherImpl struct {
}

func (h *BarcodeCRC32HasherImpl) Hash(input string) string {
	return string(crc32.Checksum([]byte(input), crc32.IEEETable))
}

func (h *BarcodeCRC32HasherImpl) CheckHash(input, hash string) bool {
	hashFromInput := h.Hash(input)
	return hashFromInput == input
}
