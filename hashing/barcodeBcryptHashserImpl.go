package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

type BarcodeBcryptHashserImpl struct {
}

func (h *BarcodeBcryptHashserImpl) Hash(input string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(input), 1)
	return string(bytes)
}

func (h *BarcodeBcryptHashserImpl) CheckHash(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}
