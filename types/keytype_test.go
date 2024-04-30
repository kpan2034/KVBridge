package types_test

import (
	. "KVBridge/types"
	"bytes"
	"testing"
)

func TestKeyType_Encode(t *testing.T) {
	key := []byte("hello_world")
	kt := NewKeyType(key)
	kt.SetHash(uint32(1))
	encoded := kt.Encode()
	expected := append([]byte{0, 0, 0, 1}, key...)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("got: %v, expected: %v", encoded, expected)
	}
}

func TestKeyType_Decode(t *testing.T) {
	key := []byte("hello_world")
	keyHash := []byte{0, 0, 0, 7}
	encoded := append(keyHash, key...)

	kt, err := DecodeToKeyType(encoded)
	if err != nil {
		t.Errorf("decode failed: %v", err)
	}
	expected := NewKeyType(key)

	if !bytes.Equal(key, kt.Key()) {
		t.Errorf("got: %v, expected: %v", encoded, expected)
	}
}
