package state

import (
	"KVBridge/config"
	. "KVBridge/types"
	"KVBridge/utils"
	"sort"
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
	N                 int
	ReplicationFactor int
	KeyRanges         []NodeRange
}

type NodeInfo struct {
	Address                string
	Status                 StatusType
	LastHeartBeatTimestamp int64
}

// GetInitialState Returns a new State struct with initial values
func GetInitialState(conf *config.Config) *State {
	id, _ := getNodeIds([]string{conf.Grpc_address}) // length of this array will always be 1
	currNodeId := id[0]
	ids, idMap := getNodeIds(conf.BootstrapServers)
	N := len(conf.BootstrapServers)

	keyRanges := make([]NodeRange, conf.ReplicationFactor)
	currNodeIdx := 0
	for ids[currNodeIdx] != currNodeId {
		currNodeIdx += 1
	}
	for i := 0; i < conf.ReplicationFactor; i++ {
		endRange := ids[(currNodeIdx-i)%N]
		startRange := ids[(currNodeIdx-i-1)%N]

		keyRanges[i] = NodeRange{StartHash: startRange, EndHash: endRange}
	}

	return &State{
		NodeType:          NodeSecondary,
		ID:                currNodeId,
		ClusterIDs:        ids,
		IDMap:             idMap,
		N:                 N,
		ReplicationFactor: conf.ReplicationFactor,
		KeyRanges:         keyRanges,
	}
}

func getNodeIds(serverAddresses []string) ([]NodeID, map[NodeID]*NodeInfo) {
	//hashGenerator := utils.SHA256HashGenerator{}
	hashGenerator := utils.Murmur3HashGenerator{}
	output := make([]NodeID, len(serverAddresses))
	idMap := make(map[NodeID]*NodeInfo)
	for i, addr := range serverAddresses {
		hash := hashGenerator.GenerateHash([]byte(addr))
		output[i] = hash

		// Add server address to map of node info
		idMap[hash] = &NodeInfo{
			Address: addr,
			Status:  StatusUP,
		}
	}

	// Sort the node ids, makes partitioner lookups simpler
	sort.Slice(output, func(i, j int) bool {
		return uint32(output[i]) < uint32(output[j])
	})

	return output, idMap
}

// helper function to easily access node info
func (s *State) GetNodeInfo(id NodeID) *NodeInfo {
	info := s.IDMap[id]
	return info
}
