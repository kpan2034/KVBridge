package state

import (
	. "KVBridge/types"
	"math/rand"
)

const (
	minNodeID int32 = 0
	maxNodeID       = 1000
)

type State struct {
	// what kind of node is this?
	NodeType
	// ID of this node
	ID NodeID
	// IDs of all nodes in the cluster
	ClusterIDs []NodeID
	// Map of nodes to any node info associated with it
	IDMap map[NodeID]*NodeInfo
	// Number of nodes in the cluster
	N int
}

type NodeInfo struct {
	Address string
}

// Returns a new State struct with initial values
func GetInititalState(sz int) *State {
	id := getNewID()
	return &State{
		NodeType:   NodeSecondary,
		ID:         id,
		ClusterIDs: []NodeID{id},
		IDMap:      make(map[NodeID]*NodeInfo),
		N:          sz,
	}
}

func getNewID() NodeID {
	return NodeID(rand.Int31n(maxNodeID-minNodeID+1) - minNodeID)
}
