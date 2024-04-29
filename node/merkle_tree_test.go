package node_test

import (
	"KVBridge/node"
	"KVBridge/types"
	"testing"
)

type TestIterator struct {
	Idx      int
	TestKeys [][]byte
}

func (it *TestIterator) Valid() bool {
	//if it.Idx >= len(it.TestKeys) {
	//	log.Printf("lol")
	//}
	return it.Idx < len(it.TestKeys)
}

func (it *TestIterator) Key() []byte {
	return it.TestKeys[it.Idx]
}

func (it *TestIterator) Next() bool {
	it.Idx += 1
	return it.Idx < len(it.TestKeys)
}

func (it *TestIterator) First() bool {
	it.Idx = 0
	return it.Idx < len(it.TestKeys)
}

func TestBuildMerkleTree(t *testing.T) {
	iter := TestIterator{Idx: 0}
	iter.TestKeys = make([][]byte, 4)
	iter.TestKeys[0] = []byte{0, 0, 0, 1}
	iter.TestKeys[1] = []byte{0, 0, 0, 2}
	iter.TestKeys[2] = []byte{0, 0, 0, 4}
	iter.TestKeys[3] = []byte{0, 0, 0, 7}
	nr := types.NodeRange{StartHash: 0, EndHash: 15}

	mt, err := node.BuildMerkleTree(nr, &iter)
	t.Logf("%v", mt.Data)
	if err != nil {
		t.Errorf("Failed building merkel tree: %s", err)
	}
	if mt.Depth != 4 {
		t.Errorf("Expected depth: 2, Actual depth %d", mt.Depth)
	}

}
