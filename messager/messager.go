package messager

import (
	"KVBridge/environment"
	"KVBridge/log"
	pb "KVBridge/proto/compiled/proto-ping"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// go generate
type Messager struct {
	*environment.Environment
	pb.UnimplementedPingServiceServer
	client pb.PingServiceClient
	conn   *grpc.ClientConn
}

func NewMessager(env *environment.Environment) *Messager {
	l := env.Named("messager")
	env = env.WithLogger(l)
	addr := "localhost:50052"
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Fatalf("did not connect: %v\n", err)
	}
	l.Debugf("connected to node at: %v\n", addr)
	c := pb.NewPingServiceClient(conn)
	m := &Messager{
		Environment: env,
		client:      c,
		conn:        conn,
	}
	l.Debugf("creating new messager: %+v", m)
	return m
}

func (m *Messager) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	m.Logger.Debugf("received ping: %v", in)
	return &pb.PingResponse{
		Resp: "hello",
	}, nil
}

func (m *Messager) PingStream(stream pb.PingService_PingStreamServer) error {
	m.Logger.Debugf("server: starting ping stream")
	counter := 1

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		msg := in.GetMsg()
		m.Logger.Debugf("server: recv count: %d, msg: %s", counter, msg)
		resp := "hello" + fmt.Sprintf("%d", counter)

		if err = stream.Send(&pb.PingResponse{
			Resp: resp,
		}); err != nil {
			m.Logger.Debugf("server: send error: %v", err)
			return err
		}
	}
}

func (m *Messager) RunPingStream(ctx context.Context) error {
	m.Logger.Debugf("starting client ping stream")

	// Set up context
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// Set up sending stream
	stream, err := m.client.PingStream(ctx)
	if err != nil {
		return err
	}

	// Initialie wait channel
	waitc := make(chan struct{})

	// Set up recv'ing stream
	go func() {
		counter := 0
		defer func() {
			m.Logger.Debugf("client: closing recv stream")
			close(waitc)
		}()
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				m.Logger.Errorf("client: recv error: %v", err)
				return
			}
			msg := in.GetResp()
			m.Logger.Debugf("client: recv count: %d, msg: %s", counter, msg)
			counter++
		}
	}()

	// run send for 100 iterations
	for i := 0; i < 100; i++ {
		msg := &pb.PingRequest{
			Msg: "hi" + fmt.Sprintf("%d", i),
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
		m.Logger.Debugf("client: send i: %d, msg: %s", i, msg)
	}
	stream.CloseSend()
	<-waitc
	return nil
}

func serialize(s string) {
	panic("unimplemented")
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
	pb.RegisterPingServiceServer(s, m)
	m.Logger.Debugf("ping server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}

	m.Logger.Debugf("ping server returning with no error\n")
	return nil
}

func (m *Messager) PingRequest(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	m.Logger.Debugf("sending ping request: ctx: %v in: %v", ctx, in)
	resp, err := m.client.Ping(ctx, in)
	if err != nil {
		m.Logger.Errorf("error response: %v", err)
	}
	m.Logger.Debugf("received ping response: resp: %v", resp)
	return resp, err
}

func InterceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		logMsg := append([]any{"msg", msg}, fields)
		l.Debugf("%v", logMsg)
	})
}
