package partitioner

import (
	. "KVBridge/types"
	"bytes"
	"encoding/binary"
	// "reflect"
)

// Partitioner interface
// Given a node and the cluster size, returns which node(s) the key belongs to
type Partitioner interface {
	GetPartitions(key []byte) ([]NodeID, error)
}

type SimplePartitioner struct {
	node NodeID
	size int32
}

func (p *SimplePartitioner) GetPartitions(key []byte) ([]NodeID, error) {
	var k NodeID
	// nodeIDSize := reflect.TypeOf(NodeID(0)).Size()
	err := binary.Read(bytes.NewReader(key), binary.BigEndian, &k)
	if err != nil {
		return nil, err
	}

	return []NodeID{k % NodeID(p.size)}, nil
}
