package messager

import (
	"KVBridge/proto/compiled/ping"
	"context"
	"fmt"
	"io"
)

// Handles Ping RPC
func (m *Messager) Ping(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {
	m.Logger.Debugf("received ping: %v", in)
	return &ping.PingResponse{
		Resp: "hello",
	}, nil
}

// Handles PingStream RPC
func (m *Messager) PingStream(stream ping.PingService_PingStreamServer) error {
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

		if err = stream.Send(&ping.PingResponse{
			Resp: resp,
		}); err != nil {
			m.Logger.Debugf("server: send error: %v", err)
			return err
		}
	}
}
