package node_test

import (
	"context"
	"testing"
	"time"
)

// func setupTest(t *testing.T) (teardownTest func(t *testing.T), node1, node2 *node.KVNode) {
// 	// Define configs for both nodes
// 	c1 := &config.Config{
// 		Address:          ":6379",
// 		Grpc_address:     "localhost:50051",
// 		LogPath:          "stdout",
// 		DataPath:         "./testing/storage",
// 		BootstrapServers: []string{"localhost:50051", "localhost:50052"},
// 	}
// 	c2 := &config.Config{
// 		Address:          ":6380",
// 		Grpc_address:     "localhost:50052",
// 		LogPath:          "stdout",
// 		DataPath:         "./testing/storage2",
// 		BootstrapServers: []string{"localhost:50051", "localhost:50052"},
// 	}
//
// 	// Create and launch both KVNodes on separate goroutines
// 	node1, err := node.NewKVNode(c1)
// 	if err != nil {
// 		t.Errorf("node1 failed: %v", err)
// 	}
// 	node2, err = node.NewKVNode(c2)
// 	if err != nil {
// 		t.Errorf("node2 failed: %v", err)
// 	}
//
// 	// Launch both servers
// 	go func() {
// 		t.Errorf("node1 failed: %v", node1.Start())
// 	}()
// 	go func() {
// 		t.Errorf("node2 failed: %v", node2.Start())
// 	}()
// 	// Wait for servers to come up
// 	time.Sleep(1000 * time.Millisecond)
//
// 	return func(t *testing.T) {
// 		os.RemoveAll("./testing")
// 	}, node1, node2
// }

func TestPingRequest(t *testing.T) {
	teardownTest, node1, _, _ := setupTest(t)
	defer teardownTest(t)

	// Perform a ping request from node1 to node2
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	_ = node1.PingEverybody(ctx)
	// if err != nil {
	// 	t.Errorf("request failed: %v", err)
	// }

	// Validate response
	// assert.Equal(t, resp.GetResp(), "hello", "response from node2 should be 'hello'")

	// Perform a ping request from node1 to node2
	ctx, cancelFunc = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// err = node1.RunPingStream(ctx)
	// if err != nil {
	// 	t.Errorf("request failed: %v", err)
	// }
}

// func TestPingStreamRequest(t *testing.T) {
// 	teardownTest, node1, _ := setupTest(t)
// 	defer teardownTest(t)
//
// }
