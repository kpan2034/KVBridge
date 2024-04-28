package node

import (
	"KVBridge/config"
	"KVBridge/environment"
	"KVBridge/log"
	"KVBridge/partitioner"
	"KVBridge/state"
	"KVBridge/storage"
	"fmt"

	"github.com/tidwall/redcon"
)

// The main struct that represents a node in a KVBridge cluster
type KVNode struct {
	// contains config, logger, state, etc
	// wrapper in an environment to make it easier to pass around to other parts of the system
	*environment.Environment
	// Storage engine
	Storage storage.StorageEngine
	// messager handle inter-node communication
	*Messager
	Partitioner partitioner.Partitioner

	// other important stuff
	client_server *redcon.Server
}

// Returns a KVNode with the specified configuration
func NewKVNode(config *config.Config) (*KVNode, func(), error) {

	node := &KVNode{
		client_server: nil,
	}

	closeFunc := func() {} // to avoid issues of func() being returned nil
	// Initialize logger
	logPath := config.LogPath
	logger := log.NewZapLogger(logPath, log.DebugLogLevel).Named(fmt.Sprintf("node@%v", config.Address))
	init_state := state.GetInitialState(config)

	// Wrap all dependencies in an env
	env := environment.New(logger, config, init_state)
	node.Environment = env

	logger.Debugf("creating kvnode with config: %v", config)

	// Initialize storage engine
	storage, err := storage.NewPebbleStorageEngine(env)
	if err != nil {
		logger.Fatalf("could not init storage engine: %v", err)
		return nil, closeFunc, err
	}
	node.Storage = storage

	// Initialize messager
	// yeah this is dumb okay
	messager := NewMessager(node)
	node.Messager = messager

	// Initialize partitioner
	partitioner, err := partitioner.GetNewPartitioner(init_state)
	if err != nil {
		logger.Errorf("could not init partitioner: %v", err)
		return nil, closeFunc, err
	}
	node.Partitioner = partitioner

	// Initialize client handler here
	srv := node.getNewClientServer()
	node.client_server = srv
	// Aggregate all close functions here and return
	closeFunc = func() {
		messager.Close()

		err := storage.Close()
		if err != nil {
			node.Logger.Errorf("error closing storage: %v", err)
		}

		err = srv.Close()
		if err != nil {
			node.Logger.Errorf("error closing client server: %v", err)
		}
	}

	return node, closeFunc, nil
}
