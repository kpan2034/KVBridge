package messager

import (
	"KVBridge/proto/compiled/ping"
	. "KVBridge/types"
	"context"
	"fmt"
	"io"
	"time"
)

func (m *Messager) PingRequest(ctx context.Context, id NodeID, in *ping.PingRequest) (*ping.PingResponse, error) {
	cl := m.getClient(id)
	return cl.PingRequest(ctx, in)
}

// Initiates PingRequest
func (cl *Client) PingRequest(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {
	cl.Logger.Debugf("sending ping request: ctx: %v in: %v", ctx, in)
	resp, err := cl.Ping(ctx, in)
	if err != nil {
		cl.Logger.Errorf("error response: %v", err)
	}
	cl.Logger.Debugf("received ping response: resp: %v", resp)
	return resp, err
}

// Initiates PingStreamRequest
func (cl *Client) RunPingStream(ctx context.Context) error {
	cl.Logger.Debugf("starting client ping stream")

	// Set up context
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// Set up sending stream
	stream, err := cl.PingStream(ctx)
	if err != nil {
		return err
	}

	// Initialie wait channel
	waitc := make(chan struct{})

	// Set up recv'ing stream
	go func() {
		counter := 0
		defer func() {
			cl.Logger.Debugf("client: closing recv stream")
			close(waitc)
		}()
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				cl.Logger.Errorf("client: recv error: %v", err)
				return
			}
			msg := in.GetResp()
			cl.Logger.Debugf("client: recv count: %d, msg: %s", counter, msg)
			counter++
		}
	}()

	// run send for 100 iterations
	for i := 0; i < 100; i++ {
		msg := &ping.PingRequest{
			Msg: "hi" + fmt.Sprintf("%d", i),
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
		cl.Logger.Debugf("client: send i: %d, msg: %s", i, msg)
	}
	stream.CloseSend()
	<-waitc
	return nil
}
