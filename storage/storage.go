package storage

import (
	. "KVBridge/environment"
	"KVBridge/types"
	"github.com/cockroachdb/pebble"
	"math"
)

var ErrNotFound = pebble.ErrNotFound

// Storage Engine Interface to support
type StorageEngine interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Close() error
	Snapshot(keyLowerBound types.NodeID, keyUpperBound types.NodeID) ([][]byte, [][]byte, error)
	GetSnapshotIter(keyLowerBound types.NodeID, keyUpperBound types.NodeID) ([]StorageIterator, error) // TODO: Decouple from pebble iterator
}

type PebbleStorageManager struct {
	*Environment
	db *pebble.DB
}

// Get the wrapper around pebble
func NewPebbleStorageEngine(env *Environment) (StorageEngine, error) {
	newEnv := env.WithLogger(env.Named("storage"))
	pb := &PebbleStorageManager{newEnv, nil}

	err := pb.Init()
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func (pb *PebbleStorageManager) Init() (err error) {
	storagePath := pb.Config.DataPath
	pb.db, err = pebble.Open(storagePath, &pebble.Options{})
	if err != nil {
		return err
	}
	return nil
}

/* TODO(kpan): I'm pretty sure the number of memory allocations in the set path kill performance.
* Likely need to abstract away direct use of the underlying PebbleDB
* Instead, use some sort of intermediate buffer to minimize number of allocations
* However, note that this is a storage-level optimization that we can look into later
 */
func (pb *PebbleStorageManager) Set(key []byte, value []byte) error {
	pb.Debugf("set request: key: %v, value: %v", string(key), string(value))
	err := pb.db.Set(key, value, &pebble.WriteOptions{})
	if err != nil {
		return err
	}

	// Flush it always for now -- need to change this obv
	err = pb.db.Flush()
	if err != nil {
		return err
	}
	pb.Debugf("set response: %v", err)
	return nil
}

func (pb *PebbleStorageManager) Get(key []byte) ([]byte, error) {
	pb.Debugf("get request: key: %v", string(key))
	value, closer, err := pb.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	pb.Debugf("get response: key: %v, value: %v", string(key), string(value))
	return value, nil
}

func (pb *PebbleStorageManager) Snapshot(keyLowerBound types.NodeID, keyUpperBound types.NodeID) ([][]byte, [][]byte, error) {
	snapshotDB := pb.db.NewSnapshot()

	iterOptions := pebble.IterOptions{
		LowerBound: types.ToBytes(keyLowerBound),
		UpperBound: types.ToBytes(keyUpperBound),
	}
	iter, err := snapshotDB.NewIter(&iterOptions)

	if err != nil {
		return nil, nil, err
	}
	var keys [][]byte
	var vals [][]byte
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, key)
		vals = append(vals, val)
	}
	if err := iter.Close(); err != nil {
		return nil, nil, err
	}
	if err := snapshotDB.Close(); err != nil {
		return nil, nil, err
	}

	return keys, vals, nil
}

type StorageIterator interface {
	Valid() bool
	Key() []byte
	Next() bool
	First() bool
	Close() error
}

func (pb *PebbleStorageManager) GetSnapshotIter(keyLowerBound types.NodeID, keyUpperBound types.NodeID) ([]StorageIterator, error) {
	snapshotDB := pb.db.NewSnapshot()
	if keyLowerBound <= keyUpperBound {
		iterOptions := pebble.IterOptions{
			LowerBound: types.ToBytes(keyLowerBound),
			UpperBound: types.ToBytes(keyUpperBound),
		}
		iter, err := snapshotDB.NewIter(&iterOptions)
		return []StorageIterator{iter}, err
	} else {
		iterOptions1 := pebble.IterOptions{
			LowerBound: types.ToBytes(0),
			UpperBound: types.ToBytes(keyUpperBound),
		}
		iter1, err := snapshotDB.NewIter(&iterOptions1)
		if err != nil {
			return nil, err
		}
		iterOptions2 := pebble.IterOptions{
			LowerBound: types.ToBytes(keyLowerBound),
			UpperBound: types.ToBytes(math.MaxUint32),
		}
		iter2, err := snapshotDB.NewIter(&iterOptions2)
		return []StorageIterator{iter1, iter2}, err
	}

}

// TODO: call this when handling graceful shutdown
func (pb *PebbleStorageManager) Close() error {
	return pb.db.Close()
}
