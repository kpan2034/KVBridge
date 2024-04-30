package node

import (
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
	key := cmd.Args[1]
	value := cmd.Args[2]

	err := node.Write(key, value)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR could not set (%s = %s): %v", string(key), string(value), err))
		return
	}
	conn.WriteString("OK")
}

func (node *KVNode) getHandler(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	key := cmd.Args[1]
	value, err := node.Read(key)
	if err != nil {
		conn.WriteError(fmt.Sprintf("ERR could not get %s", string(key)))
		return
	}
	conn.WriteBulk(value)
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
