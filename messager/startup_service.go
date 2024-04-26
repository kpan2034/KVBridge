package messager

// Stuff that the node has to do on startup go here
// Off the top of my head:
// 1. Discover other node ids (it only has bootstrap addresses)

// Handles GetNodeInfo RPC
//func (m *Messager) GetNodeInfo(ctx context.Context, in *startup.GetNodeInfoRequest) (*startup.GetNodeInfoResponse, error) {
//	m.Logger.Debugf("received ping: %v", in)
//	return &startup.GetNodeInfoResponse{
//		Id: m.State.ID,
//	}, nil
//}
