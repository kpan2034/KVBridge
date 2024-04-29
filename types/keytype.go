package types

// type of keys handled by the node
type KeyType struct {
	hashKey string
	key     []byte
}

// wrap a key in a KeyType
func NewKeyType(key []byte) *KeyType {
	// hash := getHash(key)
	return &KeyType{
		hashKey: "",
		key:     key,
	}
}

// KeyType implements the Stringer interface
func (kt *KeyType) String() string {
	return string(kt.key)
}

func (kt *KeyType) Encode() []byte {
	return kt.key
}

func DecodeToKeyType(b []byte) (*KeyType, error) {
	return &KeyType{
		hashKey: "",
		key:     b,
	}, nil
}
