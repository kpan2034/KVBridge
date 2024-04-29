package node

import (
	"KVBridge/storage"
	"errors"
	"fmt"

	"github.com/tidwall/redcon"
)

func (node *KVNode) getNewClientServer() *redcon.Server {
	// Do some setup stuff here
	addr := node.Config.Address

	// TODO(kpan) Should probably make a mux, set it up and launch a separate goroutine
	mux := redcon.NewServeMux()

	mux.HandleFunc("set", node.setHandler)
	mux.HandleFunc("get", node.getHandler)
	mux.HandleFunc("ping", node.pingHandler)
	mux.HandleFunc("quit", node.quitHandler)
	mux.HandleFunc("config", node.configHandler)
	mux.HandleFunc("del", node.commonHandler)
	mux.HandleFunc("publish", node.commonHandler)
	mux.HandleFunc("subscribe", node.commonHandler)
	mux.HandleFunc("psubscribe", node.commonHandler)

	srv := redcon.NewServer(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			command := string(cmd.Args[0])
			node.Logger.Debugf("received command: %v", command)
			mux.ServeRESP(conn, cmd)
		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			node.Logger.Debugf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			node.Logger.Debugf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)

	return srv

	// // where the mux listens and serves clients
	// err := redcon.ListenAndServe(addr,
	// 	func(conn redcon.Conn, cmd redcon.Command) {
	// 		command := string(cmd.Args[0])
	// 		node.Logger.Debugf("received command: %v", command)
	// 		mux.ServeRESP(conn, cmd)
	// 	},
	// 	switch strings.ToLower(command) {
	// 	default:
	// 		node.Logger.Errorf("unknown command: %v", command)
	// 		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	// 	case "ping":
	// 		conn.WriteString("PONG")
	// 	case "quit":
	// 		conn.WriteString("OK")
	// 		conn.Close()
	// 	case "set":
	// 		node.setHandler(conn, cmd)
	// 	case "get":
	// 		node.getHandler(conn, cmd)
	// 	case "config":
	// 		conn.WriteArray(8)
	// 		conn.WriteString("\"tcp-keepalive\"")
	// 		conn.WriteString("\"0\"")
	// 		conn.WriteString("\"io-threads\"")
	// 		conn.WriteString("\"1\"")
	// 		conn.WriteString("\"save\"")
	// 		conn.WriteString("\"3600 1 300 100 60 10000\"")
	// 		conn.WriteString("\"appendonly\"")
	// 		conn.WriteString("\"no\"")
	// 	case "del", "publish", "subscribe", "psubscribe":
	// 		conn.WriteError("ERR unsupported operation:" + string(cmd.Args[0]))
	// 	}
	// },
	// 	func(conn redcon.Conn) bool {
	// 		// Use this function to accept or deny the connection.
	// 		node.Logger.Debugf("accept: %s", conn.RemoteAddr())
	// 		return true
	// 	},
	// 	func(conn redcon.Conn, err error) {
	// 		// This is called when the connection has been closed
	// 		node.Logger.Debugf("closed: %s, err: %v", conn.RemoteAddr(), err)
	// 	},
	// )
	// if err != nil {
	// 	node.Logger.Fatalf("", err)
	// }
}

// blocks
func (node *KVNode) Start() error {
	// Do some setup stuff here
	addr := node.Config.Address

	// Initialize inter-node communication here
	// TODO(kpan): add sync primitives and graceful closing of servers
	go node.Messager.Start()

	// Initalize listener
	node.Logger.Debugf("starting server at %s", addr)

	return node.client_server.ListenAndServe()
}

func (node *KVNode) setHandler(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 3 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	// TODO: PREPEND KEY WITH HASH
	key := cmd.Args[1]
	value := cmd.Args[2]
	err := node.Storage.Set(key, value)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR could not set (%s = %s)", string(key), string(value)))
	}

	// Ping other node (testing gossip)
	nacks, err := node.ReplicateWrites(key, value)
	if err != nil {
		// just error out for now, later should just log
		conn.WriteError(fmt.Sprintf("ERR could not replicate (%s = %s)", string(key), string(value)))
	}
	node.Logger.Debugf("replicated to %d other nodes", nacks)
	conn.WriteString("OK")
}

func (node *KVNode) getHandler(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	key := cmd.Args[1]
	value, err := node.Storage.Get(key)
	// if it's a not found error, we can still go get the key from other nodes
	if errors.Is(err, storage.ErrNotFound) {
		err = nil
		// conn.WriteError(fmt.Sprintf("ERR no value associated with key: %s", string(key)))
		// return
	}
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR could not get %s", string(key)))
		return
	}

	// Reconcile value of the key with other nodes
	majorityValue, err := node.ReconcileKeyValue(key, value)
	// TODO: if this errors out, check if you can serve the local value directly
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR could not reconcile value of %s", string(key)))
		return
	}

	// Ping other node (testing gossip)
	// _ = node.Messager.PingEverybody(context.Background())
	// _, err = node.Messager.PingRequest(context.Background(), &pb.PingRequest{
	// 	Msg: "hi",
	// })
	// if err != nil {
	// 	node.Logger.Errorf("err messaging node: %v", err)
	// }
	conn.WriteBulk(majorityValue)
}

func (node *KVNode) pingHandler(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

func (node *KVNode) quitHandler(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("OK")
	conn.Close()
}

func (node *KVNode) configHandler(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteArray(8)
	conn.WriteString("\"tcp-keepalive\"")
	conn.WriteString("\"0\"")
	conn.WriteString("\"io-threads\"")
	conn.WriteString("\"1\"")
	conn.WriteString("\"save\"")
	conn.WriteString("\"3600 1 300 100 60 10000\"")
	conn.WriteString("\"appendonly\"")
	conn.WriteString("\"no\"")
}

func (node *KVNode) commonHandler(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteError("ERR unsupported operation:" + string(cmd.Args[0]))
}
