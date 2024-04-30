package node

import (
	"KVBridge/proto/compiled/replication"
	"KVBridge/storage"
	. "KVBridge/types"
	"context"
	"errors"
)

func (m *Messager) ReplicateWrite(ctx context.Context, in *replication.ReplicateWriteRequest) (*replication.ReplicateWriteResponse, error) {
	m.Logger.Debugf("received replicatewrite request: %v", in)

	inKey := in.GetKey()
	inVal := in.GetValue()
	incomingID := NodeID(in.GetId())

	// decode incoming key to a KeyType
	key, err := DecodeToKeyType(inKey)
	if err != nil {
		return nil, err
	}

	// decode incoming value to a ValueType
	incomingValue, err := DecodeToValueType(inVal)
	if err != nil {
		return nil, err
	}

	myValue, err := m.node.Storage.Get(key.Key())
	if err == storage.ErrNotFound {
		err = nil
	}
	if err != nil {
		// m.Logger.Errorf("could not key: (key:%v, value:%v): %v", key, myValue, err)
		return nil, err
	}

	// convert my value to a ValueType
	myValueDecoded, err := DecodeToValueType(myValue)
	if err != nil {
		return nil, err
	}

	shouldUpdate := m.node.ShouldUpdate(myValueDecoded, incomingValue, incomingID)

	if !shouldUpdate {
		// Ideally we should tell them that their version is outdated
		// But we just return false
		// If they want a new version they can request for it later
		return &replication.ReplicateWriteResponse{
			Ok: false,
		}, nil
	}

	err = m.node.Storage.Set(key.Encode(), incomingValue.Encode())
	if err != nil {
		m.Logger.Errorf("could not replicate: (key:%v, value:%v): %v", key, incomingValue, err)
	}

	return &replication.ReplicateWriteResponse{
		Ok: true,
	}, nil
}

func (m *Messager) GetKey(ctx context.Context, in *replication.GetKeyRequest) (*replication.GetKeyResponse, error) {
	m.Logger.Debugf("received getkey request: %v", in)

	inKey := in.GetKey()
	key, err := DecodeToKeyType(inKey)
	if err != nil {
		return nil, err
	}

	value, err := m.node.Storage.Get(key.Key())
	if errors.Is(err, storage.ErrNotFound) {
		resp := &replication.GetKeyResponse{
			Ok:    false,
			Value: nil,
		}
		return resp, nil
	}
	if err != nil {
		return nil, err
	}

	resp := &replication.GetKeyResponse{
		Ok:    true,
		Value: value,
	}
	return resp, nil
}

func (node *KVNode) ShouldUpdate(myVal *ValueType, newVal *ValueType, tieBreaker NodeID) bool {
	if newVal == nil {
		return false
	}
	if myVal == nil {
		return true
	}
	res := CompareTimestamps(myVal.Version(), newVal.Version())
	if res == 0 {
		// update only if other node has higher id
		return node.ID <= tieBreaker
	}

	// update only if my version is lower
	return (res == -1)
}

func (m *Messager) GetMerkleTree(ctx context.Context, req *replication.MerkleTreeRequest) (*replication.MerkleTreeResponse, error) {
	lb := NodeID(req.KeyRangeLowerBound)
	ub := NodeID(req.KeyRangeUpperBound)
	nr := NodeRange{StartHash: lb, EndHash: ub}
	snapshotDB := m.node.Storage.GetSnapshotDB()
	snapshotIters, err := m.node.Storage.GetSnapshotIters(lb, ub, snapshotDB)
	if err != nil {
		return nil, err
	}

	merkleTree, err := BuildMerkleTree(nr, snapshotIters)
	if err != nil {
		return nil, err
	}
	for _, ssIter := range snapshotIters {
		err := ssIter.Close()
		if err != nil {
			return nil, err
		}
	}
	err = snapshotDB.Close()
	if err != nil {
		return nil, err
	}

	merkleTreeBytes := SerializeMerkleTree(merkleTree)

	resp := replication.MerkleTreeResponse{Data: merkleTreeBytes}
	return &resp, nil
}
