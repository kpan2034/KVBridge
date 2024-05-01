package node_test

import (
	. "KVBridge/node"
	"KVBridge/proto/compiled/startup"
	. "KVBridge/types"
	"context"
	"testing"
	"time"
)

func TestMessager_StartupService(t *testing.T) {
	//teardownTest, node1, _, _, _, _, _ := setupTest2(t, 3)
	teardownTest, node1, _, _, _, _, _ := SetupTestEmulation(t)
	defer teardownTest(t)

	ids := node1.ClusterIDs

	for _, id := range ids {
		if id == node1.ID {
			continue
		}

		// GetNodeInfo
		req := &startup.GetNodeInfoRequest{}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		resp, err := node1.GetNodeInfoRequest(ctx, id, req)
		if err != nil {
			t.Errorf("error making request: %v", err)
		}

		if NodeID(resp.GetId()) != id {
			t.Errorf("want: %v, got: %v", id, resp.GetId())
		}
	}
}
