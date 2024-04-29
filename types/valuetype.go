package types

import (
	"fmt"
)

// type of the values handled by the node
type ValueType struct {
	value   []byte
	version Timestamp
}

func NewValueType(value []byte) *ValueType {
	return &ValueType{
		value:   value,
		version: NewTimestamp(),
	}
}

// get version
func (vt *ValueType) Version() Timestamp {
	return vt.version
}

// increase the local copy of the timestamp
func (vt *ValueType) Tick() {
	vt.version = vt.version.Send()
}

// update the copy of the timestamp with an incoming timestamp
func (vt *ValueType) Update(t Timestamp) {
	vt.version = vt.version.Receive(t)
}

// ValueType implements the Stringer interface
func (vt *ValueType) String() string {
	return fmt.Sprintf("%s %s", string(vt.value), vt.version)
}

func (vt *ValueType) Encode() []byte {
	// return append(vt.value, vt.version.Encode()...)
	return append(vt.version.Encode(), vt.value...)
}

func DecodeToValueType(b []byte) (*ValueType, error) {
	version, err := DecodeToTimestamp(b[:8])
	if err != nil {
		return nil, err
	}
	return &ValueType{
		value:   b[8:],
		version: version,
	}, nil
}