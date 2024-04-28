package node

import (
	"KVBridge/proto/compiled/replication"
	"KVBridge/storage"
	"context"
	"errors"
)

func (m *Messager) ReplicateWrite(ctx context.Context, in *replication.ReplicateWriteRequest) (*replication.ReplicateWriteResponse, error) {
	m.Logger.Debugf("received replicatewrite request: %v", in)

	key := in.GetKey()
	value := in.GetValue()

	err := m.node.Storage.Set(key, value)
	if err != nil {
		m.Logger.Errorf("could not replicate: (key:%v, value:%v): %v", key, value, err)
	}

	return &replication.ReplicateWriteResponse{
		Ok: true,
	}, nil
}

func (m *Messager) GetKey(ctx context.Context, in *replication.GetKeyRequest) (*replication.GetKeyResponse, error) {
	m.Logger.Debugf("received getkey request: %v", in)

	key := in.GetKey()

	value, err := m.node.Storage.Get(key)
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
