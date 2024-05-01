package node

import (
	"KVBridge/config"
	"KVBridge/environment"
	"KVBridge/log"
	"KVBridge/partitioner"
	"KVBridge/state"
	"KVBridge/storage"
	"KVBridge/types"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

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

	//
	closeFunc func()
}

func (kvNode *KVNode) getRecoverNode(keyrangeLB types.NodeID, keyrangeUB types.NodeID) (types.NodeID, error) {
	// TODO: hacky, need to connect to node to check status and if given range is served
	idx := 0
	for keyrangeUB != kvNode.Environment.State.ClusterIDs[idx] {
		idx += 1
	}
	for i := 0; i < kvNode.Environment.State.ReplicationFactor; i++ {
		candidateNodeIdx := (i + idx) % kvNode.Environment.State.N
		candidateNodeID := kvNode.Environment.State.ClusterIDs[candidateNodeIdx]
		if candidateNodeID != kvNode.ID && kvNode.Environment.State.IDMap[candidateNodeID].Status == types.StatusUP {
			return candidateNodeID, nil
		}
	}
	return 0, errors.New("No live node found")
}

// Recover Synchronizes divergent replicas (could be due to failures)
// Utilizes Anti-entropy (Merkle trees) for efficiency
func (kvNode *KVNode) Recover() error {
	kvNode.Logger.Debugf("Initiating recovery on %v", kvNode.ID)
	if kvNode.Environment.State.ReplicationFactor <= 1 {
		kvNode.Environment.Logger.Errorf("Replication Factor is %d, can not recover", kvNode.Environment.State.ReplicationFactor)
		return errors.New("Replication Factor <=1, can not recover")
	}

	// Ranges of data handled by current node
	keyRanges := kvNode.Environment.State.KeyRanges

	snapshotDB := kvNode.Storage.GetSnapshotDB()
	for _, keyrange := range keyRanges {

		// Identify which node (S) to copy over data from to current node (D)
		sourceNode, err := kvNode.getRecoverNode(keyrange.StartHash, keyrange.EndHash)
		if err != nil {
			return err
		}

		// Take a snapshot of data
		snapshotIters, err := kvNode.Storage.GetSnapshotIters(keyrange.StartHash, keyrange.EndHash, snapshotDB)
		if err != nil {
			return err
		}

		// Create merkle tree on S and D for each key range
		destMerkleTree, err := BuildMerkleTree(keyrange, snapshotIters)
		if err != nil {
			return err
		}
		for _, ssIter := range snapshotIters {
			err := ssIter.Close()
			if err != nil {
				return err
			}
		}

		srcMerkleTree, err := kvNode.FetchMerkleTree(sourceNode, uint32(keyrange.StartHash), uint32(keyrange.EndHash))
		if err != nil {
			return err
		}

		diffs, err := DiffMerkleTree(srcMerkleTree, destMerkleTree)
		if err != nil && diffs != nil {
			return err
		}

		if len(diffs) > 0 {
			err = kvNode.getClient(sourceNode).RecoverKeyRanges(context.TODO(), diffs, kvNode)
			if err != nil {
				return err
			}
		}
	}
	err := snapshotDB.Close()
	if err != nil {
		return err
	}
	kvNode.Logger.Debugf("Recovery complete on %v", kvNode.ID)
	return nil
}

// Kill a node, useful for testing purposes
func (kvNode *KVNode) Kill() error {
	//err := kvNode.storage.Close()
	//if err != nil { return err }
	os.Exit(0)
	return nil
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

	logger.Debugf("creating kvnode with config: %+v", config)

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
	closeFunc = sync.OnceFunc(
		func() {
			messager.Close()

			err := storage.Close()
			if err != nil {
				node.Logger.Errorf("error closing storage: %v", err)
			}

			err = srv.Close()
			node.Logger.Debugf("closing client handling server: %s", srv.Addr())
			if err != nil {
				node.Logger.Errorf("error closing client server: %v", err)
			}

			// Print final stats on close
			node.Logger.Infof("%s", node.Stats)
		})

	node.closeFunc = closeFunc

	return node, closeFunc, nil
}

func (node *KVNode) Close() {
	if node.closeFunc == nil {
		panic("close function not registered.")
	}
	node.closeFunc()
}
