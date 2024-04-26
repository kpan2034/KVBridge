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
	id := getNodeIds([]string{conf.Grpc_address})[0]
	return &State{
		NodeType:          NodeSecondary,
		ID:                id,
		ClusterIDs:        getNodeIds(conf.BootstrapServers),
		IDMap:             make(map[NodeID]*NodeInfo),
		N:                 len(conf.BootstrapServers),
		ReplicationFactor: conf.ReplicationFactor,
	}
}

func getNodeIds(serverAddresses []string) []NodeID {
	hashGenerator := utils.SHA256HashGenerator{}
	output := make([]string, len(serverAddresses))
	for i, addr := range serverAddresses {
		output[i] = hashGenerator.GenerateHash([]byte(addr))
	}
	// Sort the node ids, makes partitioner lookups simpler
	sort.Strings(output)
	outputNodeIds := make([]NodeID, len(serverAddresses))
	for i, hash := range output {
		outputNodeIds[i] = NodeID(hash)
	}
	return outputNodeIds
}
