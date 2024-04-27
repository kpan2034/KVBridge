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
}

type NodeInfo struct {
	Address string
}

// GetInitialState Returns a new State struct with initial values
func GetInitialState(conf *config.Config) *State {
	id, _ := getNodeIds([]string{conf.Grpc_address}) // length of this array will always be 1
	ids, idMap := getNodeIds(conf.BootstrapServers)

	return &State{
		NodeType:          NodeSecondary,
		ID:                id[0],
		ClusterIDs:        ids,
		IDMap:             idMap,
		N:                 len(conf.BootstrapServers),
		ReplicationFactor: conf.ReplicationFactor,
	}
}

func getNodeIds(serverAddresses []string) ([]NodeID, map[NodeID]*NodeInfo) {
	hashGenerator := utils.SHA256HashGenerator{}
	output := make([]string, len(serverAddresses))
	idMap := make(map[NodeID]*NodeInfo)
	for i, addr := range serverAddresses {
		hash := hashGenerator.GenerateHash([]byte(addr))
		output[i] = hash

		id := NodeID(hash)
		// Add server address to map of node info
		idMap[id] = &NodeInfo{
			Address: addr,
		}
	}

	// Sort the node ids, makes partitioner lookups simpler
	sort.Strings(output)
	outputNodeIds := make([]NodeID, len(serverAddresses))
	for i, hash := range output {
		outputNodeIds[i] = NodeID(hash)
	}
	return outputNodeIds, idMap
}

// helper function to easily access node info
func (s *State) GetNodeInfo(id NodeID) *NodeInfo {
	info := s.IDMap[id]
	return info
}
