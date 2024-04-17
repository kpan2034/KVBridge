package node_test

import (
	"KVBridge/config"
	"KVBridge/node"
	pb "KVBridge/proto/compiled/proto-ping"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNodeCommunication(t *testing.T) {
	// Define configs for both nodes
	c1 := &config.Config{
		Address:          ":6379",
		Grpc_address:     "localhost:50051",
		LogPath:          "./testing/log",
		DataPath:         "./testing/storage",
		BootstrapServers: []string{"localhost:50051", "localhost:50052"},
	}
	c2 := &config.Config{
		Address:          ":6380",
		Grpc_address:     "localhost:50052",
		LogPath:          "./testing/log2",
		DataPath:         "./testing/storage2",
		BootstrapServers: []string{"localhost:50051", "localhost:50052"},
	}

	// Create and launch both KVNodes on separate goroutines
	node1 := node.NewKVNode(c1)
	node2 := node.NewKVNode(c2)

	go func() {
		t.Errorf("node1 failed: %v", node1.Start())
	}()
	go func() {
		t.Errorf("node2 failed: %v", node2.Start())
	}()

	// Defer cleanup
	cleanupFunc := func() {
		os.RemoveAll("./testing")
	}
	defer cleanupFunc()

	// Wait for servers to come up
	time.Sleep(time.Second)

	// Perform a ping request from node1 to node2
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	resp, err := node1.PingRequest(ctx, &pb.PingRequest{
		Msg: "hi",
	})
	if err != nil {
		t.Errorf("request failed: %v", err)
	}

	// Validate response
	assert.Equal(t, resp.GetResp(), "hello", "response from node2 should be 'hello'")
}
