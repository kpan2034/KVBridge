package messager

import (
	"KVBridge/proto/compiled/startup"
	"context"
)

// Initiates PingRequest
func (cl *Client) GetNodeInfoRequest(ctx context.Context, in *startup.GetNodeInfoRequest) (*startup.GetNodeInfoResponse, error) {
	resp, err := cl.GetNodeInfo(ctx, in)
	if err != nil {
		cl.Logger.Errorf("error response: %v", err)
	}
	return resp, err
}
