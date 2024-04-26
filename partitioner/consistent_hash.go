package partitioner

import (
	"KVBridge/state"
	. "KVBridge/types"
	"KVBridge/utils"
)

type ConsistentHashPartitioner struct {
	nodeState *state.State
}

func (p *ConsistentHashPartitioner) GetPartitions(key []byte) ([]NodeID, error) {
	hashGenerator := utils.SHA256HashGenerator{}
	keyHash := hashGenerator.GenerateHash(key)

	// Use last few bits of hash (to avoid integer overflows) mod cluster_size to decide which server to use
	idx := 0
	for idx < p.nodeState.N && string(p.nodeState.ClusterIDs[idx]) < keyHash {
		idx += 1
	}
	preferenceList := make([]NodeID, p.nodeState.ReplicationFactor)
	for i := 0; i < p.nodeState.ReplicationFactor; i++ {
		preferenceList[i] = p.nodeState.ClusterIDs[(idx+i)%p.nodeState.N]
	}
	return preferenceList, nil
}
