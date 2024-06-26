package types

import (
	"encoding/binary"
)

// what kind of node is this?
type NodeType int

// type of the underlying node
type NodeID uint32

type StatusType int

type NodeRange struct {
	StartHash NodeID
	EndHash   NodeID
}

func ToBytes(nodeID NodeID) []byte {
	nodeIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(nodeIDBytes, uint32(nodeID))
	return nodeIDBytes
}

// TODO: add go generate
const (
	NodePrimary   NodeType = iota // primary
	NodeSecondary                 // secondary
)

const (
	StatusUP StatusType = iota
	StatusDOWN
)

type OpPreference int

const (
	OpLocal OpPreference = iota
	OpMajority
	OpAll
)

type OpPreferenceString string

const (
	OpLocalStr    OpPreferenceString = "local"
	OpMajorityStr                    = "majority"
	OpAllStr                         = "all"
)

func (nr NodeRange) InRange(hashedKey uint32) bool {
	return (uint32(nr.StartHash) <= hashedKey && hashedKey <= uint32(nr.EndHash)) ||
		(hashedKey <= uint32(nr.StartHash) || uint32(nr.EndHash) <= hashedKey)
}
