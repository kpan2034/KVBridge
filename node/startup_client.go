package node

import (
	"KVBridge/proto/compiled/startup"
	. "KVBridge/types"
	"context"
)

// Initiates PingRequest
func (m *Messager) GetNodeInfoRequest(ctx context.Context, id NodeID, in *startup.GetNodeInfoRequest) (*startup.GetNodeInfoResponse, error) {
	cl := m.getClient(id)
	return cl.GetNodeInfoRequest(ctx, in)
}

func (cl *Client) GetNodeInfoRequest(ctx context.Context, in *startup.GetNodeInfoRequest) (*startup.GetNodeInfoResponse, error) {
	resp, err := cl.GetNodeInfo(ctx, in)
	if err != nil {
		cl.Logger.Errorf("error response: %v", err)
	}
	return resp, err
}
