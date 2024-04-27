package messager

import (
	"KVBridge/proto/compiled/startup"
	"context"
)

// Stuff that the node has to do on startup go here
// Off the top of my head:
// 1. Discover other node ids (it only has bootstrap addresses)

// Handles GetNodeInfo RPC
func (s *startupServer) GetNodeInfo(ctx context.Context, in *startup.GetNodeInfoRequest) (*startup.GetNodeInfoResponse, error) {
	s.Logger.Debugf("received getnodeinfo request: %v", in)
	return &startup.GetNodeInfoResponse{
		Id: string(s.State.ID),
	}, nil
}
