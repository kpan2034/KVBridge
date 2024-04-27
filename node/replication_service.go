package node

import (
	"KVBridge/proto/compiled/replication"
	"context"
)

func (m *Messager) ReplicateWrite(ctx context.Context, in *replication.ReplicateWriteRequest) (*replication.ReplicateWriteResponse, error) {
	m.Logger.Debugf("received replicatewrite request: %v", in)

	key := in.GetKey()
	value := in.GetValue()

	err := m.node.storage.Set(key, value)
	if err != nil {
		m.Logger.Errorf("could not replicate: (key:%v, value:%v): %v", key, value, err)
	}

	return &replication.ReplicateWriteResponse{
		Ok: true,
	}, nil
}
