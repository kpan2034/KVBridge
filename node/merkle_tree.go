package node

import (
	"KVBridge/storage"
	"KVBridge/types"
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

func BuildMerkleTree(nr types.NodeRange, iters []storage.StorageIterator) (*MerkleTree, error) {

	depth := min(MAX_DEPTH, int(math.Ceil(math.Log2(float64(nr.EndHash-nr.StartHash)))))
	data := make([]uint32, int(math.Pow(2, float64(depth)+1))-1)
	if len(iters) == 1 {
		iter := iters[0]
		iter.First()
		_ = buildTreeUtil(0, 0, depth, nr.StartHash, nr.EndHash, &data, iter)
	} else if len(iters) == 2 {
		iters[0].First()
		_ = buildTreeUtil(0, 0, depth, 0, nr.EndHash, &data, iters[0])
		iters[1].First()
		_ = buildTreeUtil(0, 0, depth, nr.StartHash, math.MaxUint32, &data, iters[1])
	} else {
		log.Fatalf("Length of iters is expected to be 1 or 2, received %d", len(iters))
	}

	mt := MerkleTree{Depth: depth, RangeLowerBound: uint32(nr.StartHash), RangeUpperBound: uint32(nr.EndHash), Data: data}
	return &mt, nil
}

func buildTreeUtil(loc int, currDepth int, maxDepth int, lb types.NodeID, ub types.NodeID, data *[]uint32, iter storage.StorageIterator) uint32 {
	if currDepth == maxDepth {
		// Leaf nodes
		acc := (*data)[loc]
		exit := false
		for iter.Valid() && !exit {
			keytype, _ := types.DecodeToKeyType(iter.Key())
			if keytype.Hash() > uint32(ub) {
				exit = true
			} else if keytype.Hash() <= uint32(ub) && keytype.Hash() >= uint32(lb) {
				acc = acc ^ keytype.Hash()
				iter.Next()
			} else {
				iter.Next()
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

	//if s.RangeLowerBound <= s.RangeUpperBound {
	//	return diffUtil(0, s.RangeLowerBound, s.RangeUpperBound, s, d), nil
	//} else {
	//	return append(
	//		diffUtil(0, 0, s.RangeLowerBound, s, d),
	//		diffUtil(0, s.RangeUpperBound, math.MaxUint32, s, d)...), nil
	//}
	return diffUtil(0, 0, math.MaxUint32, s, d), nil
}

func diffUtil(loc int, lb uint32, ub uint32, s *MerkleTree, d *MerkleTree) []types.NodeRange {
	if ub < lb {
		log.Fatalf("diffUtil expects ub>=lb")
	}

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
