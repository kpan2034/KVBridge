package types

import (
	"encoding/binary"
)

// type of keys handled by the node
type KeyType struct {
	hash uint32
	key  []byte
}

// wrap a key in a KeyType
func NewKeyType(key []byte) *KeyType {
	hashGenerator := Murmur3HashGenerator{}
	hash := uint32(hashGenerator.GenerateHash(key))
	return &KeyType{
		hash: hash,
		key:  key,
	}
}

func (kt *KeyType) Key() []byte {
	return kt.key
}
func (kt *KeyType) Hash() uint32 {
	return kt.hash
}

// KeyType implements the Stringer interface
func (kt *KeyType) String() string {
	return string(kt.key)
}

func (kt *KeyType) Encode() []byte {
	hashBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(hashBytes, kt.hash)
	return append(hashBytes, kt.key...)
}

func DecodeToKeyType(b []byte) (*KeyType, error) {
	return &KeyType{
		hash: binary.BigEndian.Uint32(b[:4]),
		key:  b[4:],
	}, nil
}
