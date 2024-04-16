package node

type nodeType int

const (
	nodePrimary nodeType = iota
	nodeSecondary
)

type State struct {
	nodeType nodeType
}
