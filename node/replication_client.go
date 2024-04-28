package node

import (
	"KVBridge/proto/compiled/replication"
	// "sync"
	// . "KVBridge/types"
	"context"
	"sync/atomic"
)

// Initiates ReplicateWriteRequest
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
		resp, err := cl.ReplicateWrite(context.TODO(), in)
		if err != nil {
			return 0, err
		}
		m.Logger.Debugf("replicated to %s: ok: ", id, resp.GetOk)

		if resp.GetOk() {
			acks.Add(1)
		}
		// }()
	}

	return int(acks.Load()), nil
}

// Gets value of a keys from all nodes, and returns the majority value.
// If no majority value exists then returns votedValue is nil
// If the node calling ReconcileKeyValue does not have a value associated with the provided key
// then myValue can be nil
func (m *Messager) ReconcileKeyValue(key []byte, myValue []byte) (votedValue []byte, err error) {
	m.Logger.Debugf("getting value of key: %s from all nodes", key)

	// convert []byte key to a string -- there needs to be a better way of doing this btw
	voteMap := make(map[string]int)
	if myValue != nil {
		voteMap[string(myValue)] = 1
	}

	for _, id := range m.ClusterIDs {
		if id == m.ID {
			continue
		}

		// Make the request
		in := &replication.GetKeyRequest{
			Key: key,
		}

		cl := m.getClient(id)
		resp, err := cl.GetKey(context.TODO(), in)
		if err != nil {
			return nil, err
		}
		ok := resp.GetOk()
		value := resp.GetValue()
		m.Logger.Debugf("node %s (ok:%v) returned (key: %v, value: %v)", id, ok, key, value)

		// Add returned value to the vote map
		voteMap[string(value)]++
	}

	// Extract majority value from the map and return
	value, hasMajority := getMajorityKey(m.N, voteMap)
	if !hasMajority {
		return nil, nil
	}

	// TODO: if majority value is different from myValue, start a background routine to write the new value

	return []byte(value), nil
}

// ideally you write this as a generic but ok
func getMajorityKey(N int, voteMap map[string]int) (string, bool) {
	theta := N/2 + 1
	for k, v := range voteMap {
		if v >= theta {
			return k, true
		}
	}

	return "", false
}
