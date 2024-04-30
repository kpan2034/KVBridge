package node

import (
	"KVBridge/storage"
	"KVBridge/types"
	"encoding/binary"
	"log"
	"math"
)

const MAX_DEPTH int = 15

type MerkleTree struct {
	Depth           int
	RangeLowerBound uint32
	RangeUpperBound uint32
	Data            []uint32
}

func BuildMerkleTree(nr types.NodeRange, iter storage.StorageIterator) (*MerkleTree, error) {

	depth := min(MAX_DEPTH, int(math.Ceil(math.Log2(float64(nr.EndHash-nr.StartHash+1)))))
	data := make([]uint32, int(math.Pow(2, float64(depth)+1))-1)
	iter.First()
	_ = buildTreeUtil(0, 0, depth, nr.StartHash, nr.EndHash, &data, iter)
	mt := MerkleTree{Depth: depth, RangeLowerBound: uint32(nr.StartHash), RangeUpperBound: uint32(nr.EndHash), Data: data}
	return &mt, nil
}

func buildTreeUtil(loc int, currDepth int, maxDepth int, lb types.NodeID, ub types.NodeID, data *[]uint32, iter storage.StorageIterator) uint32 {
	if currDepth == maxDepth {
		// Leaf nodes
		acc := uint32(0)
		exit := false
		for iter.Valid() && !exit {
			keyHashBytes := iter.Key()[:4]
			k := binary.BigEndian.Uint32(keyHashBytes)
			if k > uint32(ub) {
				exit = true
			} else if k <= uint32(ub) && k >= uint32(lb) {
				acc = acc ^ k
				iter.Next()
			} else {
				log.Fatalf("Unexpected key %d received in buildTreeUtil", k)
			}
		}
		(*data)[loc] = acc
	} else {
		midPoint := lb + (ub-lb)/2
		leftChildHash := buildTreeUtil(2*loc+1, currDepth+1, maxDepth, lb, midPoint, data, iter)
		rightChildHash := buildTreeUtil(2*loc+2, currDepth+1, maxDepth, midPoint+1, ub, data, iter)
		(*data)[loc] = leftChildHash ^ rightChildHash
	}
	return (*data)[loc]
}

func DiffMerkleTree(s *MerkleTree, d *MerkleTree) ([]types.NodeRange, error) {
	if s.Depth != d.Depth {
		log.Fatalf("Unexpected call to DiffMerkleTree with source tree depth: %d, dest tree depth: %d ", s.Depth, d.Depth)
	}
	if s.RangeLowerBound != d.RangeLowerBound || s.RangeUpperBound != d.RangeUpperBound {
		log.Fatalf("Unexpected call to DiffMerkleTree with source tree range [%d, %d], dest tree range [%d, %d]",
			s.RangeLowerBound, s.RangeUpperBound, d.RangeLowerBound, d.RangeUpperBound)
	}

	return diffUtil(0, s.RangeLowerBound, s.RangeUpperBound, s, d), nil
}

func diffUtil(loc int, lb uint32, ub uint32, s *MerkleTree, d *MerkleTree) []types.NodeRange {
	if loc >= len(s.Data) || s.Data[loc] == d.Data[loc] {
		return []types.NodeRange{}
	} else {
		if loc >= int(math.Pow(2, float64(s.Depth))-1) {
			// leaf nodes
			nr := types.NodeRange{
				StartHash: types.NodeID(lb),
				EndHash:   types.NodeID(ub),
			}
			return []types.NodeRange{nr}
		} else {
			midPoint := lb + (ub-lb)/2
			leftDiffs := diffUtil(2*loc+1, lb, midPoint, s, d)
			rightDiffs := diffUtil(2*loc+2, midPoint+1, ub, s, d)
			return append(leftDiffs, rightDiffs...)
		}
	}
}

func SerializeMerkleTree(mt *MerkleTree) []uint32 {
	metadata := []uint32{uint32(mt.Depth), mt.RangeLowerBound, mt.RangeUpperBound}
	return append(metadata, mt.Data...)
}

func DeserializeMerkleTree(data []uint32) (*MerkleTree, error) {
	depth := data[0]
	lb := data[1]
	ub := data[2]
	treeData := data[3:]
	return &MerkleTree{
		Depth:           int(depth),
		RangeLowerBound: lb,
		RangeUpperBound: ub,
		Data:            treeData}, nil
}
