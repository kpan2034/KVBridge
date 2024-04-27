package node

import (
	"KVBridge/proto/compiled/replication"
	// "sync"
	// . "KVBridge/types"
	"context"
	"sync/atomic"
)

// Initiates PingRequest
func (m *Messager) ReplicateWrites(key, value []byte) (nacks int, err error) {
	m.Logger.Debugf("replicating write: (%d, %d)", key, value)

	var acks atomic.Int32
	// var wg sync.WaitGroup
	// wg.Add(len(m.ClusterIDs))

	for _, id := range m.ClusterIDs {
		if id == m.ID {
			continue
		}
		// go func() {
		in := &replication.ReplicateWriteRequest{
			Key:   key,
			Value: value,
		}
		cl := m.getClient(id)
		resp, err := cl.ReplicateWriteRequest(context.TODO(), in)
		if err != nil {
			return 0, err
		}
		m.Logger.Debugf("replicated to %v: ok: ", id, resp.GetOk)

		if resp.GetOk() {
			acks.Add(1)
		}
		// }()
	}

	return int(acks.Load()), nil
}

func (cl *Client) ReplicateWriteRequest(ctx context.Context, in *replication.ReplicateWriteRequest) (*replication.ReplicateWriteResponse, error) {
	resp, err := cl.ReplicateWrite(ctx, in)
	if err != nil {
		cl.Logger.Errorf("error response: %v", err)
	}
	return resp, err
}
