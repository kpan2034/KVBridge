package storage

import (
	"KVBridge/log"
	"github.com/cockroachdb/pebble"
)

// Storage Engine Interface to support
type StorageEngine interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
}

type PebbleStorageManager struct {
	db     *pebble.DB
	logger log.Logger
}

// Get the wrapper around pebble
func NewPebbleStorageEngine() (StorageEngine, error) {
	pb := &PebbleStorageManager{}

	err := pb.Init()
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func (pb *PebbleStorageManager) Init() (err error) {
	pb.db, err = pebble.Open("./tmp/storage", &pebble.Options{})
	if err != nil {
		return err
	}
	return nil
}

func (pb *PebbleStorageManager) Set(key []byte, value []byte) error {
	err := pb.db.Set(key, value, &pebble.WriteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (pb *PebbleStorageManager) Get(key []byte) ([]byte, error) {
	value, _, err := pb.db.Get(key)
	// defer closer.Close()
	if err != nil {
		return nil, err
	}
	return value, nil
}
