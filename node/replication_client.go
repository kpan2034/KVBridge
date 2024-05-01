package node

import (
	"KVBridge/proto/compiled/replication"
	. "KVBridge/types"
	"context"
	"errors"
	"io"
)

var ErrReplicateWriteFailed = errors.New("replicate write failed")
var ErrReconcileReadFailed = errors.New("reconcile read failed")
var ErrTimeout = errors.New("timeout")

// Initiates ReplicateWriteRequest
func (m *Messager) ReplicateWrites(key *KeyType, value *ValueType) (nacks int, err error) {
	m.Logger.Debugf("replicating write: (%v, %v)", key, value)

	acks := make(chan struct{}, m.N)    // struct{} is apparently a null kinda value
	errChan := make(chan struct{}, m.N) // This channel buffer size may not be ideal for large clusters

	ctx, cancel := context.WithTimeout(context.TODO(), m.Config.Timeout)
	defer cancel()

	partitions, _ := m.node.Partitioner.GetPartitions(key.Key()) // this can never return an error
	for _, id := range partitions {
		m.Logger.Debugf("replicating to node:%v, timeout: %s", id, m.Config.Timeout)
		go func(id NodeID) {
			if id == m.ID {
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

	for {
		select {
		case <-acks:
			num_acks++
			m.Logger.Debugf("ack on ackChan, num_acks: %d", num_acks)
			if num_acks >= req_acks {
				return num_acks, nil
			}

		case <-errChan:
			num_errs++
			m.Logger.Debugf("err on errChan, num_errs: %d, num_acks: %d", num_errs, num_acks)
			// if too many errors, then fail
			if num_errs > (m.N - req_acks) {
				return num_acks, ErrReplicateWriteFailed
			}

		case <-ctx.Done():
			m.Logger.Debugf("context done, errs: %d, acks: %d", num_errs, num_acks)
			return num_acks, ErrTimeout
		}
	}
}

// Gets value of a keys from all nodes, and returns the majority value.
// If no majority value exists then returns votedValue is nil
// If the node calling ReconcileKeyValue does not have a value associated with the provided key
// then myValue can be nil
func (m *Messager) ReconcileKeyValue(key *KeyType, myValue *ValueType) (votedValue []byte, err error) {
	m.Logger.Debugf("reconciling value of key: %s", key)

	// note: for comparision, []byte keys are converted to a string
	// sigh -- there needs to be a better way of doing this btw
	voteMap := make(map[string]int)
	if myValue != nil && myValue.Value() != nil {
		voteMap[string(myValue.Value())] = 1
	}

	type valWithId struct {
		vt *ValueType
		id NodeID
	}
	valChan := make(chan valWithId, m.N)
	errChan := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.TODO(), m.Config.Timeout)
	defer cancel()

	for _, id := range m.ClusterIDs {
		if id == m.ID {
			valChan <- valWithId{myValue, id}
			continue
		}

		go func(id NodeID) {
			// Make the request
			in := &replication.GetKeyRequest{
				Key: key.Encode(),
			}

			cl := m.getClient(id)
			resp, err := cl.GetKey(ctx, in)
			if err != nil {
				m.Logger.Errorf("error getting value of key from node %v: %v", id, err)
				errChan <- struct{}{}
				return
			}
			ok := resp.GetOk()
			value := resp.GetValue()

			// decode value to value type
			vt, err := DecodeToValueType(value)
			if err != nil {
				m.Logger.Errorf("error decoding value of key from node %v: %v", id, err)
				errChan <- struct{}{}
				return
			}
			m.Logger.Debugf("node %v (ok:%v) returned (key: %v, value: %v)", id, ok, key, vt)
			valChan <- valWithId{vt, id}
		}(id)
	}

	req_acks := m.GetReadThreshold()
	num_acks := 0
	num_errs := 0 // you can calc this using resps and acks btw

	var valueToReturn ValueType
	for {
		select {
		case val := <-valChan:
			num_acks++
			m.Logger.Debugf("got ack: val to return: %#v, val:%#v, num_acks: %d", &valueToReturn, val.vt, num_acks)
			if m.node.ShouldUpdate(&valueToReturn, val.vt, val.id) {
				m.Logger.Debugf("updating val to return")
				valueToReturn = *val.vt
			}
			if num_acks >= req_acks {
				return valueToReturn.Value(), nil
			}
		case <-errChan:
			num_errs++
			// if too many errors, then fail
			if num_errs > (m.N - req_acks) {
				return nil, ErrReconcileReadFailed
			}
		case <-ctx.Done():
			m.Logger.Debugf("context done, errs: %d, acks: %d", num_errs, num_acks)
			return nil, ErrTimeout
		}
	}

	// Add returned value to the vote map
	// We ignore versions when voting

}

func (m *Messager) FetchMerkleTree(nodeID NodeID, lb uint32, ub uint32) (*MerkleTree, error) {
	merkleTreeReq := replication.MerkleTreeRequest{KeyRangeLowerBound: lb, KeyRangeUpperBound: ub}
	m.Logger.Debugf("%v", m.client_map)
	merkleTreeResp, err := m.getClient(nodeID).GetMerkleTree(context.TODO(), &merkleTreeReq)

	if err != nil {
		return nil, err
	}
	tree, err := DeserializeMerkleTree(merkleTreeResp.Data)
	return tree, err
}

func (c *Client) RecoverKeyRanges(ctx context.Context, keyRanges []NodeRange, n *KVNode) error {

	protoKeyRanges := make([]*replication.GetKeysInRangesRequest_KeyRange, len(keyRanges))
	for i, keyRange := range keyRanges {
		protoKeyRanges[i] = &replication.GetKeysInRangesRequest_KeyRange{
			KeyRangeLowerBound: uint32(keyRange.StartHash),
			KeyRangeUpperBound: uint32(keyRange.EndHash),
		}
	}

	req := replication.GetKeysInRangesRequest{KeyRangeList: protoKeyRanges}
	stream, err := c.GetKeysInRanges(ctx, &req)
	if err != nil {
		return nil
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ok := resp.GetOk()
		key := resp.GetKey()
		val := resp.GetValue()
		if ok {
			err = n.WriteWithReplicate(key, val, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
