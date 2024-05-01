package node

import (
	"KVBridge/storage"
	. "KVBridge/types"
	"errors"
)

func (node *KVNode) Read(key []byte) ([]byte, error) {
	kt := NewKeyType(key)

	isOwner := node.ownsKey(kt.Hash())

	// forward request to some owneer
	if !isOwner {
		return nil, errors.New("not the owner")
		// return forwardRead
	}

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
	// val, err := node.ReconcileKeyValue(kt, vt)
	resChan := make(chan []byte, 1)
	errChan := make(chan error, 1)
	go node.ReconcileKeyValueAsync(kt, vt, resChan, errChan)
	select {
	case val := <-resChan:
		return val, nil
	case err := <-errChan:
		return nil, err
	}
	// if err != nil {
	// 	return nil, err
	// }
	//
	// // extract the actual value from majorityValue
	// // majorityValue, err := DecodeToValueType(val)
	// // if err != nil {
	// // 	return nil, err
	// // }
	//
	// return val, nil
}
