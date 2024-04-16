package node

import (
	"KVBridge/log"
	"KVBridge/storage"
	"fmt"
	"strings"

	"github.com/tidwall/redcon"
)

// The main struct that represents a node in a KVBridge cluster
type KVNode struct {
	config  *KVNodeConfig
	address string
	logger  log.Logger
	storage storage.StorageEngine
}

func (node *KVNode) Start() error {
	// Do some setup stuff here
	addr := node.config.addr

	// Defer cleanup stuff here
	defer node.storage.Close()
	// Initalize listener
	node.logger.Debugf("started server at %s", addr)
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
