package node

import (
	"KVBridge/log"
	"KVBridge/storage"
)

type KVNodeConfig struct {
	addr    string
	logPath string
}

func DefaultConfig() *KVNodeConfig {
	return &KVNodeConfig{
		addr:    ":6379",
		logPath: "./tmp/log",
	}
}

// Returns a KVNode with the specified configuration
func (config *KVNodeConfig) Build() *KVNode {
	// Address at which node exists
	addr := config.addr

	// Initialize logger
	logPath := config.logPath
	logger := log.NewZapLogger(logPath, log.DebugLogLevel)

	// Initialize storage engine
	storage, err := storage.NewPebbleStorageEngine()
	if err != nil {
		logger.Fatalf("could not init storage engine: %v", err)
	}

	// Return the node
	node := &KVNode{
		config:  config,
		address: addr,
		logger:  logger,
		storage: storage,
	}

	return node
}
