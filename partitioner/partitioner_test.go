package partitioner_test

import (
	"KVBridge/config"
	"KVBridge/node"
	"os"
	"testing"
	"time"
)

func setupTest(t *testing.T) (teardownTest func(t *testing.T), node1, node2, node3 *node.KVNode) {
	os.RemoveAll("./testing")
	// Define configs for both nodes
	c1 := &config.Config{
		Address:          ":6379",
		Grpc_address:     "localhost:50051",
		LogPath:          "stdout",
		DataPath:         "./testing/storage",
		BootstrapServers: []string{"localhost:50051", "localhost:50052", "localhost:50053"},
	}
	c2 := &config.Config{
		Address:          ":6380",
		Grpc_address:     "localhost:50052",
		LogPath:          "stdout",
		DataPath:         "./testing/storage2",
		BootstrapServers: []string{"localhost:50051", "localhost:50052", "localhost:50053"},
	}
	c3 := &config.Config{
		Address:          ":6381",
		Grpc_address:     "localhost:50053",
		LogPath:          "stdout",
		DataPath:         "./testing/storage3",
		BootstrapServers: []string{"localhost:50051", "localhost:50052", "localhost:50053"},
	}

	// Create and launch both KVNodes on separate goroutines
	var err error
	node1, err = node.NewKVNode(c1)
	if err != nil {
		t.Errorf("node1 creation failed: %v", err)
	}
	node2, err = node.NewKVNode(c2)
	if err != nil {
		t.Errorf("node2 creation failed: %v", err)
	}
	node3, err = node.NewKVNode(c3)
	if err != nil {
		t.Errorf("node3 creation failed: %v", err)
	}

	// Launch both servers
	go func() {
		t.Errorf("node1 failed: %v", node1.Start())
	}()
	go func() {
		t.Errorf("node2 failed: %v", node2.Start())
	}()
	go func() {
		t.Errorf("node3 failed: %v", node3.Start())
	}()

	// Wait for servers to come up
	time.Sleep(1000 * time.Millisecond)

	return func(t *testing.T) {
		os.RemoveAll("./testing")
	}, node1, node2, node3
}

func TestSimplePartitioner_GetPartitions(t *testing.T) {
	teardownTest, node1, _, _ := setupTest(t)
	defer teardownTest(t)

	partitions, err := node1.Partitioner.GetPartitions([]byte("Hello, world"))
	//fmt.Println(partitions)
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}
	p := partitions[0]
	if !(p == node1.ClusterIDs[0] || p == node1.ClusterIDs[1] || p == node1.ClusterIDs[2]) {
		t.Errorf("Unknown NodeID assiged: %s", p)
	}
}
