package types_test

import (
	. "KVBridge/types"
	"bytes"
	"testing"
)

func TestValueType_Encode(t *testing.T) {
	val := []byte("hello_world")
	randomUnixTime := int64(0x0001000200030004)
	ts := FromUnixTime(randomUnixTime).Tick().Tick().Tick().Tick()

	vt := NewValueTypeWithTimestamp(val, ts)
	encoded := vt.Encode()

	expectedEncodedTS := []byte{0x04, 0x00, 0x03, 0x00, 0x02, 0x00, 0x01, 0x00}
	expectedEncodedValue := []byte("hello_world")
	expected := append(expectedEncodedTS, expectedEncodedValue...)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("got: %v, expected: %v", encoded, expected)
	}
}

func TestValueType_Decode(t *testing.T) {
	encodedTS := []byte{0x04, 0x00, 0x03, 0x00, 0x02, 0x00, 0x01, 0x00}
	encodedValue := []byte("hello_world")
	encodedVT := append(encodedTS, encodedValue...)

	decoded, err := DecodeToValueType(encodedVT)
	if err != nil {
		t.Errorf("could not decode vt: %v", err)
	}

	val := []byte("hello_world")
	randomUnixTime := uint64(0x0001000200030004)

	if !bytes.Equal(val, decoded.Value()) {
		t.Errorf("got: %v, expected: %v", decoded.Value(), val)
	}

	if randomUnixTime != decoded.Version().Uint64() {
		t.Errorf("got: %v, expected: %v", decoded.Version().Uint64(), randomUnixTime)
	}

}
