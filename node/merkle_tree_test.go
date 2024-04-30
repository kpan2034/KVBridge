package node_test

import (
	"KVBridge/node"
	"KVBridge/storage"
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

func (it *TestIterator) Value() []byte {
	return []byte{}
}

func (it *TestIterator) Next() bool {
	it.Idx += 1
	return it.Idx < len(it.TestKeys)
}

func (it *TestIterator) First() bool {
	it.Idx = 0
	return it.Idx < len(it.TestKeys)
}

func (it *TestIterator) Close() error { return nil }

func TestBuildMerkleTree(t *testing.T) {
	iter := TestIterator{Idx: 0}
	iter.TestKeys = make([][]byte, 4)
	iter.TestKeys[0] = []byte{0, 0, 0, 1}
	iter.TestKeys[1] = []byte{0, 0, 0, 2}
	iter.TestKeys[2] = []byte{0, 0, 0, 4}
	iter.TestKeys[3] = []byte{0, 0, 0, 8}
	nr := types.NodeRange{StartHash: 0, EndHash: 15}

	mt, err := node.BuildMerkleTree(nr, []storage.StorageIterator{&iter})
	//t.Logf("%v", mt.Data)
	if err != nil {
		t.Errorf("Failed building merkel tree: %s", err)
	}
	if mt.Depth != 15 {
		t.Errorf("Expected depth: 2, Actual depth %d", mt.Depth)
	}
	dataSum := int64(0)
	for _, num := range mt.Data {
		dataSum += int64(num)
	}
	if dataSum == 0 {
		t.Errorf("No data inserted into tree!")
	}
}

func TestDiffMerkleTree(t *testing.T) {
	iter1 := TestIterator{Idx: 0}
	iter1.TestKeys = make([][]byte, 3)
	iter1.TestKeys[0] = []byte{0, 0, 0, 1}
	iter1.TestKeys[1] = []byte{0, 0, 0, 2}
	iter1.TestKeys[2] = []byte{0, 0, 0, 8}
	nr := types.NodeRange{StartHash: 0, EndHash: 15}
	destMerkleTree, err := node.BuildMerkleTree(nr, []storage.StorageIterator{&iter1})
	if err != nil {
		t.Errorf("Failed building merkel tree: %s", err)
	}
	iter2 := TestIterator{Idx: 0}
	iter2.TestKeys = make([][]byte, 4)
	iter2.TestKeys[0] = []byte{0, 0, 0, 1}
	iter2.TestKeys[1] = []byte{0, 0, 0, 2}
	iter2.TestKeys[2] = []byte{0, 0, 0, 4}
	iter2.TestKeys[3] = []byte{0, 0, 0, 8}
	sourceMerkleTree, err := node.BuildMerkleTree(nr, []storage.StorageIterator{&iter2})
	if err != nil {
		t.Errorf("Failed building merkel tree: %s", err)
	}

	diffs, err := node.DiffMerkleTree(sourceMerkleTree, destMerkleTree)
	if err != nil {
		t.Errorf("Failed DiffMerkleTree merkel trees: %s", err)
	}

	if len(diffs) != 1 {
		t.Errorf("Expected len of diffs 1, actual %d", len(diffs))
	}
}

func TestDiffMerkleTree2(t *testing.T) {
	iter1 := TestIterator{Idx: 0}
	iter1.TestKeys = make([][]byte, 3)
	iter1.TestKeys[0] = []byte{0, 0, 0, 1}
	iter1.TestKeys[1] = []byte{0, 0, 0, 2}
	iter1.TestKeys[2] = []byte{0, 0, 0, 8}
	nr := types.NodeRange{StartHash: 16, EndHash: 1}
	destMerkleTree, err := node.BuildMerkleTree(nr, []storage.StorageIterator{&iter1, &TestIterator{}})
	if err != nil {
		t.Errorf("Failed building merkel tree: %s", err)
	}
	iter2 := TestIterator{Idx: 0}
	iter2.TestKeys = make([][]byte, 4)
	iter2.TestKeys[0] = []byte{0, 0, 0, 1}
	iter2.TestKeys[1] = []byte{0, 0, 0, 2}
	iter2.TestKeys[2] = []byte{0, 0, 0, 4}
	iter2.TestKeys[3] = []byte{0, 0, 0, 8}
	sourceMerkleTree, err := node.BuildMerkleTree(nr, []storage.StorageIterator{&iter2, &TestIterator{}})
	if err != nil {
		t.Errorf("Failed building merkel tree: %s", err)
	}

	diffs, err := node.DiffMerkleTree(sourceMerkleTree, destMerkleTree)
	if err != nil {
		t.Errorf("Failed DiffMerkleTree merkel trees: %s", err)
	}

	// diffs is still 1, node range filtering only at iterator level, design can be improved
	if len(diffs) != 1 {
		t.Errorf("Expected len of diffs 1, actual %d", len(diffs))
	}
}
