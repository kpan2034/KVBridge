package node

import (
	"KVBridge/environment"
	"KVBridge/log"
	"KVBridge/proto/compiled/ping"
	"KVBridge/proto/compiled/replication"
	"KVBridge/proto/compiled/startup"
	. "KVBridge/types"
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

type Messager struct {
	*environment.Environment
	pingServer
	startupServer
	replicationServer
	client_map map[NodeID]*Client
	clients    []Client
	node       *KVNode // backpointer to kv node, to be used internally by the messager struct only
	server     *grpc.Server
	// need to find a better way to doing this. really messed up the design of this.
}

type pingServer struct {
	*environment.Environment
	ping.UnimplementedPingServiceServer
}

type startupServer struct {
	*environment.Environment
	startup.UnimplementedStartupServiceServer
}

type replicationServer struct {
	*environment.Environment
	replication.UnimplementedReplicationServiceServer
}

// The messager has a client that is associated with each node
// On startup, it talks to other nodes and maps them to a node id
type Client struct {
	*environment.Environment
	*grpc.ClientConn
	ping.PingServiceClient
	startup.StartupServiceClient
	replication.ReplicationServiceClient
}

func NewMessager(node *KVNode) *Messager {
	env := node.Environment
	l := env.Named("messager")
	env = env.WithLogger(l)
	clients := make([]Client, 0, len(env.ClusterIDs))
	client_map := make(map[NodeID]*Client)

	// create a new client for each server in the config
	for _, id := range env.ClusterIDs {
		if id == env.ID {
			continue
		}
		l.Debugf("creating client for: %v", id)

		addr := env.GetNodeInfo(id).Address

		// Set up a connection to that node.
		conn, err := grpc.Dial(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(logging.UnaryClientInterceptor(InterceptorLogger(l))),
		)
		if err != nil {
			l.Fatalf("did not connect: %v\n", err)
		}
		l.Debugf("connected to node at: %v\n", addr)

		new_client := newClient(env.WithLogger(env.Named(addr)), conn)
		clients = append(clients, new_client)
		client_map[id] = &new_client
	}

	m := &Messager{
		Environment:       env,
		pingServer:        pingServer{},
		startupServer:     startupServer{},
		replicationServer: replicationServer{},
		client_map:        client_map,
		clients:           clients,
		node:              node,
		server:            nil,
	}

	l.Debugf("creating new messager: %+v", m)
	return m
}

func newClient(env *environment.Environment, conn *grpc.ClientConn) Client {
	return Client{
		env,
		conn,
		ping.NewPingServiceClient(conn),
		startup.NewStartupServiceClient(conn),
		replication.NewReplicationServiceClient(conn),
	}
}

func (m *Messager) Start() error {
	lis, err := net.Listen("tcp", m.Config.Grpc_address)
	if err != nil {
		return err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(m.Logger)),
		),
	)
	m.server = s

	var wg sync.WaitGroup

	numServers := 0

	numServers++
	pinger := &pingServer{
		Environment: m.Environment.WithLogger(m.Logger.Named("ping_service")),
	}
	ping.RegisterPingServiceServer(s, pinger)

	numServers++
	starter := &startupServer{
		Environment: m.Environment.WithLogger(m.Logger.Named("startup_service")),
	}
	startup.RegisterStartupServiceServer(s, starter)

	numServers++
	// replicator := &replicationServer{
	// 	Environment: m.Environment.WithLogger(m.Logger.Named("replication_service")),
	// }
	replication.RegisterReplicationServiceServer(s, m)

	wg.Add(numServers)
	go func() {
		defer wg.Done()
		pinger.Logger.Debugf("ping server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			pinger.Logger.Errorf("ping server error: %v", err)
		}
		pinger.Logger.Debugf("ping server returning with no error\n")
	}()

	go func() {
		defer wg.Done()
		starter.Logger.Debugf("startup server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			starter.Logger.Errorf("startup server error: %v", err)
		}
		starter.Logger.Debugf("startup server returning with no error\n")
	}()

	go func() {
		defer wg.Done()
		m.Logger.Debugf("replication server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			m.Logger.Errorf("replication server error: %v", err)
		}
		m.Logger.Debugf("replication server returning with no error\n")
	}()

	wg.Wait()

	return nil
}

func InterceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		logMsg := append([]any{"msg", msg}, fields)
		l.Debugf("%v", logMsg)
	})
}

func (m *Messager) getClient(id NodeID) *Client {
	client := m.client_map[id]
	return client
}

func (m *Messager) Close() {
	m.server.GracefulStop()
}
