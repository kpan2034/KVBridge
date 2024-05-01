package state

import (
	"KVBridge/config"
	. "KVBridge/types"
	"fmt"
	"log"
	"sort"
	"sync/atomic"
	"time"
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
	*Stats
}

type Stats struct {
	NumReads     atomic.Uint64
	NumWrites    atomic.Uint64
	ErrReads     atomic.Uint64
	ErrWrites    atomic.Uint64
	NumReplicate atomic.Uint64
	NumOutdated  atomic.Uint64
	TotalDelay   atomic.Uint64
}

func (s *Stats) CountRead() {
	s.NumReads.Add(1)
}

func (s *Stats) CountWrite() {
	s.NumWrites.Add(1)
}

func (s *Stats) CountErrRead() {
	s.ErrReads.Add(1)
}

func (s *Stats) CountErrWrite() {
	s.ErrWrites.Add(1)
}

func (s *Stats) CountReplicate() {
	s.NumReplicate.Add(1)
}

func (s *Stats) CountInconsistent() {
	s.NumOutdated.Add(1)
}

func (s *Stats) AddDelay(d uint64) {
	s.TotalDelay.Add(d)
}

func (s *Stats) MeasureDelay(vt ValueType) {
	diff := int64(NewTimestamp().Uint64()) - int64(vt.Version().Uint64())
	if diff >= 0 {
		s.AddDelay(uint64(diff))
		s.CountReplicate()
	}
}

func (s *Stats) PrintStatHelper() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Printf("%s", s)
		}
	}
}

func (s *Stats) String() string {
	d := time.Duration(s.TotalDelay.Load() / s.NumReplicate.Load())
	str := fmt.Sprintf("reads:%d\terr_reads:%d\twrites:%d\terr_writes:%d\taverage_lag:%s\treplicate_requests:%d",
		s.NumReads.Load(),
		s.ErrReads.Load(),
		s.NumWrites.Load(),
		s.ErrWrites.Load(),
		d,
		s.NumReplicate.Load())

	return str
}

type NodeInfo struct {
	Address                string
	Status                 StatusType
	LastHeartBeatTimestamp int64
}

// GetInitialState Returns a new State struct with initial values
func GetInitialState(conf *config.Config) *State {
	// length of this array will always be 1
	id, _ := getNodeIds([]string{conf.Grpc_address})
	currNodeId := id[0]
	ids, idMap := getNodeIds(conf.BootstrapServers)
	N := len(conf.BootstrapServers)

	keyRanges := make([]NodeRange, conf.ReplicationFactor)
	currNodeIdx := 0
	for ids[currNodeIdx] != currNodeId {
		currNodeIdx += 1
	}
	for i := 0; i < conf.ReplicationFactor; i++ {
		endRange := ids[(currNodeIdx-i+N)%N]
		startRange := ids[(currNodeIdx-i+N-1)%N]

		keyRanges[i] = NodeRange{StartHash: startRange, EndHash: endRange}
	}

	stats := &Stats{}
	stats.NumReplicate.Store(1)

	return &State{
		NodeType:          NodeSecondary,
		ID:                currNodeId,
		ClusterIDs:        ids,
		IDMap:             idMap,
		N:                 N,
		ReplicationFactor: conf.ReplicationFactor,
		KeyRanges:         keyRanges,
		Stats:             stats,
	}
}

func getNodeIds(serverAddresses []string) ([]NodeID, map[NodeID]*NodeInfo) {
	//hashGenerator := utils.SHA256HashGenerator{}
	hashGenerator := Murmur3HashGenerator{}
	output := make([]NodeID, len(serverAddresses))
	idMap := make(map[NodeID]*NodeInfo)
	for i, addr := range serverAddresses {
		hash := hashGenerator.GenerateHash([]byte{addr[len(addr)-1]})
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
