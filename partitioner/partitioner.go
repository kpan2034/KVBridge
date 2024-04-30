package partitioner

import (
	"KVBridge/state"
	. "KVBridge/types"
	"strconv"
)

// Partitioner interface
type Partitioner interface {
	// GetPartitions returns which node(s) the key belongs to given a key
	GetPartitions(key []byte) ([]NodeID, error)
}

func GetNewPartitioner(s *state.State) (Partitioner, error) {
	// TODO: make partitioner choice an argument
	//partitioner := SimplePartitioner{s}
	partitioner := ConsistentHashPartitioner{s}
	return &partitioner, nil
}

type SimplePartitioner struct {
	nodeState *state.State
}

// Simple partitioner, works only if no nodes enter/exit the system
func (p *SimplePartitioner) GetPartitions(key []byte) ([]NodeID, error) {
	hashGenerator := SHA256HashGenerator{}
	var keyHash = hashGenerator.GenerateHash(key)
	var cluster_size int = p.nodeState.N
	// Use last few bits of hash (to avoid integer overflows) mod cluster_size to decide which server to use
	decodedInt, err := strconv.ParseInt(keyHash[len(keyHash)-10:len(keyHash)], 16, 64)
	if err != nil {
		return nil, err
	}
	return []NodeID{p.nodeState.ClusterIDs[decodedInt%int64(cluster_size)]}, nil
}
