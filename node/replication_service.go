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

	key := in.GetKey()
	val := in.GetValue()
	incomingID := NodeID(in.GetId())

	// decode incoming kt to a KeyType
	// kt, err := DecodeToKeyType(key)
	// if err != nil {
	// 	return nil, err
	// }

	// decode incoming value to a ValueType
	incomingVt, err := DecodeToValueType(val)
	if err != nil {
		return nil, err
	}
	// measure delay here
	m.node.State.MeasureDelay(*incomingVt)

	// shouldUpdate := true
	myValue, err := m.node.Storage.Get(key)
	if err == storage.ErrNotFound {
		m.Logger.Debugf("key not found: %v", key)
		err = nil
		// shouldUpdate = true
	}
	if err != nil {
		m.Logger.Errorf("could not key: (key:%v, value:%v): %v", key, myValue, err)
		return nil, err
	}

	// otherwise, you got a valid value, so compare to new value
	// convert my value to a ValueType, if value found
	myValueVt, err := DecodeToValueType(myValue)
	if err != nil {
		return nil, err
	}
	shouldUpdate := m.node.ShouldUpdate(myValueVt, incomingVt, incomingID)
	if !shouldUpdate {
		// Respond with our version
		return &replication.ReplicateWriteResponse{
			Resp: &replication.ReplicateWriteResponse_Value{Value: myValue},
		}, nil
	}

	err = m.node.Storage.Set(key, val)
	if err != nil {
		m.Logger.Errorf("could not replicate: (key:%v, value:%v): %v", key, val, err)
	}

	return &replication.ReplicateWriteResponse{
		Resp: &replication.ReplicateWriteResponse_Ok{Ok: true},
	}, nil
}

func (m *Messager) GetKey(ctx context.Context, in *replication.GetKeyRequest) (*replication.GetKeyResponse, error) {
	m.Logger.Debugf("received getkey request: %v", in)

	inKey := in.GetKey()
	key, err := DecodeToKeyType(inKey)
	if err != nil {
		return nil, err
	}

	value, err := m.node.Storage.Get(inKey)
	if errors.Is(err, storage.ErrNotFound) {
		m.Logger.Debugf("not found: key: %v, returning false to requestor", key)
		resp := &replication.GetKeyResponse{
			Ok:    false,
			Value: nil,
		}
		return resp, nil
	}
	if err != nil {
		m.Logger.Errorf("error getting key %v", inKey)
		return nil, err
	}

	m.Logger.Debugf("found key: %v value: %v, returning true to requestor", inKey, value)
	resp := &replication.GetKeyResponse{
		Ok:    true,
		Value: value,
	}
	return resp, nil
}

func (node *KVNode) ShouldUpdate(myVal *ValueType, newVal *ValueType, tieBreaker NodeID) bool {
	node.Logger.Debugf("checking if should update: mine: %#v, theirs: %#v, tie: %v", myVal, newVal, tieBreaker)
	if newVal == nil || newVal.Value() == nil {
		return false
	}
	if myVal == nil || myVal.Value() == nil {
		return true
	}
	res := CompareTimestamps(myVal.Version(), newVal.Version())
	node.Logger.Debugf("timestamps compared, res: %d", res)
	if res == 0 {
		// update only if other node has higher id
		node.Logger.Debugf("same timestamp, comparing node ids: %s", node.ID <= tieBreaker)
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

func (m *Messager) GetKeysInRanges(req *replication.GetKeysInRangesRequest, stream replication.ReplicationService_GetKeysInRangesServer) error {
	snapshotDB := m.node.Storage.GetSnapshotDB()
	for _, keyRange := range req.GetKeyRangeList() {
		lb := keyRange.GetKeyRangeLowerBound()
		ub := keyRange.GetKeyRangeUpperBound()

		ssIters, err := m.node.Storage.GetSnapshotIters(NodeID(lb), NodeID(ub), snapshotDB)
		if err != nil {
			return err
		}

		for _, iter := range ssIters {
			iter.First()
			for iter.Valid() {
				key := iter.Key()
				val := iter.Value()
				resp := replication.GetKeysInRangesResponse{Ok: true, Key: key, Value: val}
				err := stream.Send(&resp)
				if err != nil {
					return err
				}
				iter.Next()
			}
			err := iter.Close()
			if err != nil {
				return err
			}
		}

	}
	err := snapshotDB.Close()
	if err != nil {
		return err
	}
	return nil
}
