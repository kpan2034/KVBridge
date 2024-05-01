package node_test

import (
	"KVBridge/node"
	. "KVBridge/types"
	"bytes"
	"testing"
)

func TestMessager_ReplicateWrites(t *testing.T) {
	teardownTest, node1, node2, node3, _, _, _ := setupTest2(t, 3)
	defer teardownTest(t)

	keys := []string{"1", "2", "5", "9", "10"}
	values := []string{"1", "2", "5", "9", "10"}

	ids := node1.ClusterIDs

	for _, id := range ids {
		if id == node1.ID {
			continue
		}

		for i := range keys {
			key := []byte(keys[i])
			value := []byte(values[i])

			kt := NewKeyType(key)
			vt := NewValueType(value)

			// Write key value pair to local storage
			node1.Storage.Set(kt.Encode(), vt.Encode())
			// Replicate (key, value) to other nodes
			nacks, err := node1.ReplicateWrites(kt, vt)
			if err != nil {
				t.Errorf("error making request: %v", err)
			}
			if nacks != len(ids)-1 {
				t.Errorf("did not receive all acks: expected: %d, got :%d", len(ids)-1, nacks)
			}

			// Check that replication happened succesfully
			for _, n := range []*node.KVNode{node2, node3} {
				val, err := n.Storage.Get(kt.Encode())
				if err != nil {
					t.Errorf("node %v: could not get key %v", n.ID, key)
				}
				localVt, err := DecodeToValueType(val)
				if err != nil {
					t.Errorf("node %v: could not decode value %v", n.ID, val)
				}
				if !bytes.Equal(value, localVt.Value()) {
					t.Errorf("node %v: expected: %v, got: %v", n.ID, value, val)
				}
			}
		}
	}
}

func TestMessager_ReconcileKey(t *testing.T) {
	teardownTest, node1, node2, node3, _, _, _ := setupTest2(t, 3)
	defer teardownTest(t)

	keys := []string{"1", "2", "5", "9", "10"}
	values := []string{"1", "2", "5", "9", "10"}
	latestValue := []byte("XXXX")

	ids := node1.ClusterIDs
	latestVt := NewValueType(latestValue)

	for _, id := range ids {
		if id == node1.ID {
			continue
		}

		for i := range keys {
			key := []byte(keys[i])
			value := []byte(values[i])
			kt := NewKeyType(key)

			// Write correct key value pairs to local storage of other nodes
			err := node2.WriteWithReplicate(key, value, false)
			if err != nil {
				t.Errorf("error writing to node %v: (key: %v, value: %v)", node2.ID, key, value)
			}
			err = node2.WriteWithReplicate(key, value, false)
			if err != nil {
				t.Errorf("error writing to node %v: (key: %v, value: %v)", node3.ID, key, value)
			}

			// Write latest key value pair to local storage
			err = node1.WriteWithReplicate(key, latestValue, false)
			if err != nil {
				t.Errorf("error writing to node %v: (key: %v, value: %v)", node1.ID, key, latestValue)
			}

			// reconcile value from other nodes
			reconciledValue, err := node1.ReconcileKeyValue(kt, latestVt)
			if err != nil {
				t.Errorf("error making request: %v", err)
			}

			if !bytes.Equal(reconciledValue, latestValue) {
				t.Errorf("node %v: expected reconciled value: %v, got: %v", node1.ID, reconciledValue, latestValue)
			}
		}
	}
}
