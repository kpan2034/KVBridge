package node

import (
	"KVBridge/storage"
	. "KVBridge/types"
	"errors"
)

func (node *KVNode) Write(key []byte, value []byte) error {

	kt := NewKeyType(key)

	// Get the old version of the key
	old_value, err := node.Storage.Get(kt.Encode())
	if errors.Is(err, storage.ErrNotFound) {
		err = nil
	}
	if err != nil {
		return err
	}

	vt, err := DecodeToValueType(old_value)
	if err != nil {
		return err
	}

	// update the value and tick the version
	vt.UpdateValue(value)
	vt.Tick()

	// store new version + value
	err = node.Storage.Set(kt.Encode(), vt.Value())
	if err != nil {
		return err
	}
	node.Logger.Debugf("wrote (%v:%v) to local storage", kt, vt)

	// replicate write to other node
	nacks, err := node.ReplicateWrites(kt, vt)
	if err != nil {
		// just error out for now, later should just log
		return err
	}
	node.Logger.Debugf("replicated to %d other nodes", nacks)

	return nil
}
