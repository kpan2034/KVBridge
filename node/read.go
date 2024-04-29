package node

import (
	"KVBridge/storage"
	"errors"
)

func (node *KVNode) Read(key []byte) ([]byte, error) {
	value, err := node.Storage.Get(key)
	// if it's a not found error, we can still go get the key from other nodes
	if errors.Is(err, storage.ErrNotFound) {
		err = nil
	}

	majorityValue, err := node.ReconcileKeyValue(key, value)
	// TODO: if this errors out, check if you can serve the local value directly
	if err != nil {
		return nil, err
	}

	return majorityValue, nil
}
