package node

import (
	"KVBridge/storage"
	. "KVBridge/types"
	"errors"
)

func (node *KVNode) Read(key []byte) ([]byte, error) {
	kt := NewKeyType(key)
	value, err := node.Storage.Get(kt.Encode())
	// if it's a not found error, we can still go get the key from other nodes
	if errors.Is(err, storage.ErrNotFound) {
		err = nil
	}

	vt, err := DecodeToValueType(value)
	if err != nil {
		return nil, err
	}

	// We send the entire value over, include version information
	// TODO: do this async for speedup
	val, err := node.ReconcileKeyValue(kt, vt)
	// TODO: if this errors out, check if you can serve the local value directly
	if err != nil {
		return nil, err
	}

	// extract the actual value from majorityValue
	majorityValue, err := DecodeToValueType(val)
	if err != nil {
		return nil, err
	}

	return majorityValue.Value(), nil
}
