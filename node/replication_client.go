package node

import (
	"KVBridge/proto/compiled/replication"
	"errors"

	// "sync"
	. "KVBridge/types"
	"context"
)

var ErrReplicateWriteFailed = errors.New("replicate write failed")
var ErrTimeout = errors.New("timeout")

// Initiates ReplicateWriteRequest
func (m *Messager) ReplicateWrites(key *KeyType, value *ValueType) (nacks int, err error) {
	m.Logger.Debugf("replicating write: (%v, %v)", key, value)

	acks := make(chan struct{}, m.N)    // struct{} is apparently a null kinda value
	errChan := make(chan struct{}, m.N) // This channel buffer size may not be ideal for large clusters

	ctx, cancel := context.WithTimeout(context.TODO(), m.Config.Timeout)
	defer cancel()

	for _, id := range m.ClusterIDs {
		m.Logger.Debugf("replicating to node:%v, timeout: %s", id, m.Config.Timeout)
		go func(id NodeID) {
			if id == m.ID {
				m.Logger.Debugf("ignoring self: id:%v", id)
				acks <- struct{}{}
				return
			}

			in := &replication.ReplicateWriteRequest{
				Key:   key.Encode(),
				Value: value.Encode(),
			}
			cl := m.getClient(id)
			resp, err := cl.ReplicateWrite(ctx, in)
			if err != nil {
				m.Logger.Errorf("replicate to node:%v: failed: %v", id, err)
				errChan <- struct{}{}
				return
			}

			switch resp.GetResp().(type) {
			case *replication.ReplicateWriteResponse_Ok:
				m.Logger.Debugf("replicate to node:%v: ok", id)
				acks <- struct{}{}

			case *replication.ReplicateWriteResponse_Value:
				m.Logger.Errorf("replicate to node:%v: failed", id)
				errChan <- struct{}{}
			}
		}(id)
	}

	req_acks := m.GetWriteThreshold()
	num_acks := 0
	num_errs := 0 // you can calc this using resps and acks btw
	num_resps := 0

	for {
		select {
		case <-acks:
			num_acks++
			m.Logger.Debugf("ack on ackChan, num_acks: %d, num_resps: %d", num_acks, num_resps)
			if num_acks >= req_acks {
				return num_acks, nil
			}

		case <-errChan:
			m.Logger.Debugf("err on errChan, num_errs: %d, num_resps: %d", num_errs, num_resps)
			if num_errs >= req_acks {
				return num_acks, ErrReplicateWriteFailed
			}

		case <-ctx.Done():
			m.Logger.Debugf("context done, responses: %d, acks: %d", num_resps, num_acks)
			return num_acks, ErrTimeout
		}
		num_resps++
	}
}

// Gets value of a keys from all nodes, and returns the majority value.
// If no majority value exists then returns votedValue is nil
// If the node calling ReconcileKeyValue does not have a value associated with the provided key
// then myValue can be nil
func (m *Messager) ReconcileKeyValue(key *KeyType, myValue *ValueType) (votedValue []byte, err error) {
	m.Logger.Debugf("getting value of key: %s from all nodes", key)

	// note: for comparision, []byte keys are converted to a string
	// sigh -- there needs to be a better way of doing this btw
	voteMap := make(map[string]int)
	if myValue != nil {
		voteMap[string(myValue.Value())] = 1
	}

	for _, id := range m.ClusterIDs {
		if id == m.ID {
			continue
		}

		// Make the request
		in := &replication.GetKeyRequest{
			Key: key.Encode(),
		}

		cl := m.getClient(id)
		resp, err := cl.GetKey(context.TODO(), in)
		if err != nil {
			return nil, err
		}
		ok := resp.GetOk()
		value := resp.GetValue()

		// decode value to value type
		vt, err := DecodeToValueType(value)
		if err != nil {
			return nil, err
		}
		m.Logger.Debugf("node %v (ok:%v) returned (key: %v, value: %v)", id, ok, key, vt)

		// Add returned value to the vote map
		// We ignore versions when voting
		voteMap[string(vt.Value())]++
	}

	// Extract majority value from the map and return
	value, hasMajority := getMajorityKey(m.N, voteMap)
	m.Logger.Debugf("majority?: %s; value: %v voteMap: %#v", hasMajority, value, voteMap)
	if !hasMajority {
		return nil, nil
	}

	// TODO: if majority value is different from myValue, start a background routine to write the new value

	m.Logger.Debugf("returning: %s", value)
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
