package types_test

import (
	. "KVBridge/types"
	"bytes"
	"testing"
)

func TestKeyType_Encode(t *testing.T) {
	key := []byte("hello_world")
	kt := NewKeyType(key)
	encoded := kt.Encode()
	expected := key

	if !bytes.Equal(encoded, expected) {
		t.Errorf("got: %v, expected: %v", encoded, expected)
	}
}

func TestKeyType_Decode(t *testing.T) {
	key := []byte("hello_world")
	encoded := []byte("hello_world")

	kt, err := DecodeToKeyType(key)
	if err != nil {
		t.Errorf("decode failed: %v", err)
	}
	expected := NewKeyType(key)

	if !bytes.Equal(key, kt.Key()) {
		t.Errorf("got: %v, expected: %v", encoded, expected)
	}
}
