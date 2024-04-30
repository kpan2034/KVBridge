package node

import (
	"KVBridge/proto/compiled/replication"
	"KVBridge/types"
	"io"
	"log"
	// "sync"
	. "KVBridge/types"
	"context"
	"sync/atomic"
)

// Initiates ReplicateWriteRequest
func (m *Messager) ReplicateWrites(key *KeyType, value *ValueType) (nacks int, err error) {
	m.Logger.Debugf("replicating write: (%v, %v)", key, value)

	var acks atomic.Int32
	// var wg sync.WaitGroup
	// wg.Add(len(m.ClusterIDs))

	for _, id := range m.ClusterIDs {
		if id == m.ID {
			continue
		}
		// go func() {
		in := &replication.ReplicateWriteRequest{
			Key:   key.Encode(),
			Value: value.Encode(),
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
		m.Logger.Debugf("node %s (ok:%v) returned (key: %v, value: %v)", id, ok, key, vt)

		// Add returned value to the vote map
		// We ignore versions when voting
		voteMap[string(vt.Value())]++
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

func (m *Messager) FetchMerkleTree(nodeID types.NodeID, lb uint32, ub uint32) (*MerkleTree, error) {
	merkleTreeReq := replication.MerkleTreeRequest{KeyRangeLowerBound: lb, KeyRangeUpperBound: ub}
	log.Printf("%v", m.client_map)
	merkleTreeResp, err := m.getClient(nodeID).GetMerkleTree(context.TODO(), &merkleTreeReq)

	if err != nil {
		return nil, err
	}
	tree, err := DeserializeMerkleTree(merkleTreeResp.Data)
	return tree, err
}

func (c *Client) RecoverKeyRanges(ctx context.Context, keyRanges []types.NodeRange, n *KVNode) error {

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
