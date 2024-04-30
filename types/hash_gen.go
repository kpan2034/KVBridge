package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spaolacci/murmur3"
)

type HashGenerator interface {
	GenerateHash(plaintext []byte) NodeID
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

type Murmur3HashGenerator struct{}

func (h *Murmur3HashGenerator) GenerateHash(plaintext []byte) NodeID {
	hash := murmur3.New32()

	_, err := hash.Write(plaintext)
	if err != nil {
		// TODO: add proper logging
		fmt.Printf("ERROR: unable to generate hash %s", err)
	}

	return NodeID(hash.Sum32())
}
