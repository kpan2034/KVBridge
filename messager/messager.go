package messager

import (
	"KVBridge/environment"
	"KVBridge/log"
	"KVBridge/proto/compiled/ping"
	"KVBridge/proto/compiled/startup"
	. "KVBridge/types"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

type Messager struct {
	*environment.Environment
	pingServer
	startupServer
	client_map map[NodeID]*Client
	clients    []Client
}

type pingServer struct {
	*environment.Environment
	ping.UnimplementedPingServiceServer
}

type startupServer struct {
	*environment.Environment
	startup.UnimplementedStartupServiceServer
}

// The messager has a client that is associated with each node
// On startup, it talks to other nodes and maps them to a node id
type Client struct {
	*environment.Environment
	*grpc.ClientConn
	ping.PingServiceClient
	startup.StartupServiceClient
}

func NewMessager(env *environment.Environment) *Messager {
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
		Environment: env,
		clients:     clients,
		client_map:  client_map,
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

	var wg sync.WaitGroup

	pinger := &pingServer{
		Environment: m.Environment.WithLogger(m.Logger.Named("ping_service")),
	}

	starter := &startupServer{
		Environment: m.Environment.WithLogger(m.Logger.Named("startup_service")),
	}
	ping.RegisterPingServiceServer(s, pinger)
	startup.RegisterStartupServiceServer(s, starter)

	wg.Add(2)
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

	wg.Wait()

	return nil
}

func (m *Messager) PingEverybody(ctx context.Context) interface{} {
	// TODO: make this a buffered channel

	wg := sync.WaitGroup{}

	// Make async Ping requests
	for _, client := range m.clients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			resp, err := client.PingRequest(ctx, &ping.PingRequest{
				Msg: fmt.Sprintf("hello from: %v", m.Address),
			})
			if err != nil {
				m.Logger.Debugf("ping error: %v", err)
				return
			}
			// Ideally you send the resp over a channel so that some receiver function can handle them
			// as they come in, one at a time
			m.Logger.Debugf("recv value: %+v", resp)
		}()
	}
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
