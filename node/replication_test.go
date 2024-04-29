package node_test

import (
	"KVBridge/node"
	"bytes"
	"testing"
)

func TestMessager_ReplicateWrites(t *testing.T) {
	teardownTest, node1, node2, node3 := setupTest(t)
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

			// Write key value pair to local storage
			node1.Storage.Set(key, value)
			// Replicate (key, value) to other nodes
			nacks, err := node1.ReplicateWrites(key, value)
			if err != nil {
				t.Errorf("error making request: %v", err)
			}
			if nacks != len(ids)-1 {
				t.Errorf("did not receive all acks: expected: %d, got :%d", len(ids)-1, nacks)
			}

			// Check that replication happened succesfully
			for _, n := range []*node.KVNode{node2, node3} {
				val, err := node2.Storage.Get(key)
				if err != nil {
					t.Errorf("node %v: could not get key %v", n.ID, key)
				}
				if !bytes.Equal(value, val) {
					t.Errorf("node %v: expected: %v, got: %v", n.ID, value, val)
				}
			}
		}
	}
}

func TestMessager_ReconcileKey(t *testing.T) {
	teardownTest, node1, node2, node3 := setupTest(t)
	defer teardownTest(t)

	keys := []string{"1", "2", "5", "9", "10"}
	values := []string{"1", "2", "5", "9", "10"}
	staleValue := []byte("0")

	ids := node1.ClusterIDs

	for _, id := range ids {
		if id == node1.ID {
			continue
		}

		for i := range keys {
			key := []byte(keys[i])
			value := []byte(values[i])

			// Write invalid key value pair to local storage
			err := node1.Storage.Set(key, staleValue)
			if err != nil {
				t.Errorf("error writing to node %v: (key: %v, value: %v)", node1.ID, key, staleValue)
			}
			// Write correct key value pairs to local storage of other nodes
			err = node2.Storage.Set(key, value)
			if err != nil {
				t.Errorf("error writing to node %v: (key: %v, value: %v)", node2.ID, key, value)
			}
			err = node3.Storage.Set(key, value)
			if err != nil {
				t.Errorf("error writing to node %v: (key: %v, value: %v)", node3.ID, key, value)
			}
			// reconcile value from other nodes
			majValue, err := node1.ReconcileKeyValue(key, staleValue)
			if err != nil {
				t.Errorf("error making request: %v", err)
			}

			if !bytes.Equal(majValue, value) {
				t.Errorf("node %v: expected reconciled value: %v, got: %v", node1.ID, majValue, value)
			}
		}
	}
}
