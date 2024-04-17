package node

import (
	"KVBridge/config"
	"KVBridge/environment"
	"KVBridge/log"
	pb "KVBridge/proto/compiled/proto-ping"
	"KVBridge/storage"
	"context"
	"fmt"
	"strings"

	"github.com/tidwall/redcon"
)

// The main struct that represents a node in a KVBridge cluster
type KVNode struct {
	config  *config.Config
	address string
	logger  log.Logger
	storage storage.StorageEngine
	state   *State
	// handle inter-node communication
	messager *Messager
}

func (node *KVNode) Start() error {
	// Do some setup stuff here
	addr := node.config.Address

	// Defer cleanup stuff here
	defer node.storage.Close()

	// Initialize inter-node communication here
	// TODO(kpan): add sync primitives and graceful closing of servers
	go node.messager.Start()

	// Initalize listener
	node.logger.Debugf("started server at %s", addr)

	// TODO(kpan) Should probably make a mux, set it up and launch a separate goroutine
	// where the mux listens and serves clients
	redcon.NewServeMux()
	err := redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			command := string(cmd.Args[0])
			node.logger.Debugf("received command: %v", command)
			switch strings.ToLower(command) {
			default:
				// node.logger.Error("unknown command: %v", command)
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case "ping":
				conn.WriteString("PONG")
			case "quit":
				conn.WriteString("OK")
				conn.Close()
			case "set":
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				key := cmd.Args[1]
				value := cmd.Args[2]
				err := node.storage.Set(key, value)
				if err != nil {
					conn.WriteError(fmt.Sprintf("ERR could not set (%s = %s)", string(key), string(value)))
				}

				// Ping other node (testing gossip)
				node.messager.PingRequest(context.Background(), &pb.PingRequest{
					Msg: "hi",
				})
				conn.WriteString("OK")
			case "get":
				if len(cmd.Args) != 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				key := cmd.Args[1]
				value, err := node.storage.Get(key)
				if err != nil {
					conn.WriteError(fmt.Sprintf("ERR could not get %s", string(key)))
					return
				}
				// Ping other node (testing gossip)
				_, err = node.messager.PingRequest(context.Background(), &pb.PingRequest{
					Msg: "hi",
				})
				if err != nil {
					node.logger.Errorf("err messaging node: %v", err)
				}
				conn.WriteBulk(value)
			case "config":
				conn.WriteArray(8)
				conn.WriteString("\"tcp-keepalive\"")
				conn.WriteString("\"0\"")
				conn.WriteString("\"io-threads\"")
				conn.WriteString("\"1\"")
				conn.WriteString("\"save\"")
				conn.WriteString("\"3600 1 300 100 60 10000\"")
				conn.WriteString("\"appendonly\"")
				conn.WriteString("\"no\"")
			case "del", "publish", "subscribe", "psubscribe":
				conn.WriteError("ERR unsupported operation:" + string(cmd.Args[0]))
			}
		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			node.logger.Debugf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			node.logger.Debugf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		node.logger.Fatalf("", err)
	}

	return nil
}

// Returns a KVNode with the specified configuration
func NewKVNode(config *config.Config) *KVNode {
	addr := config.Address

	// Initialize logger
	logPath := config.LogPath
	logger := log.NewZapLogger(logPath, log.DebugLogLevel)

	// Wrap all dependencies in an env
	env := environment.New(logger, config)

	logger.Debugf("creating kvnode with config: %v", config)

	// Initialize storage engine
	storage, err := storage.NewPebbleStorageEngine(env)
	if err != nil {
		logger.Fatalf("could not init storage engine: %v", err)
	}

	// Initialize messager
	messager := NewMessager(env)

	// Return the node
	node := &KVNode{
		config:   config,
		address:  addr,
		logger:   logger,
		storage:  storage,
		messager: messager,
	}

	return node
}
