package types

// what kind of node is this?
type NodeType int

// type of the underlying node
type NodeID int32

// TODO: add go generate
const (
	NodePrimary   NodeType = iota // primary
	NodeSecondary                 // secondary
)
