package node_test

import (
	"KVBridge/node"
	"bytes"
	"testing"
)

func TestMessager_ReplicationService(t *testing.T) {
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
