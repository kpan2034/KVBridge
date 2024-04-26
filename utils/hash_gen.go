package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type HashGenerator interface {
	GenerateHash(plaintext []byte) string
}

type SHA256HashGenerator struct {
}

func (h *SHA256HashGenerator) GenerateHash(plaintext []byte) string {
	// Create a new SHA256 hash object
	hash := sha256.New()

	// Write the string data to the hash object
	_, err := hash.Write(plaintext)
	if err != nil {
		// TODO: add proper logging
		fmt.Printf("ERROR: unable to generate hash %s", err)
	}

	// Get the finalized hash result as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
