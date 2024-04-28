package node

import (
	"KVBridge/proto/compiled/ping"
	"context"
	"fmt"
	"io"
)

// Handles Ping RPC
func (ps *pingServer) Ping(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {
	ps.Logger.Debugf("received ping: %v", in)
	return &ping.PingResponse{
		Resp: "hello",
	}, nil
}

// Handles PingStream RPC
func (ps *pingServer) PingStream(stream ping.PingService_PingStreamServer) error {
	ps.Logger.Debugf("server: starting ping stream")
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
		ps.Logger.Debugf("server: recv count: %d, msg: %s", counter, msg)
		resp := "hello" + fmt.Sprintf("%d", counter)

		if err = stream.Send(&ping.PingResponse{
			Resp: resp,
		}); err != nil {
			ps.Logger.Debugf("server: send error: %v", err)
			return err
		}
	}
}
