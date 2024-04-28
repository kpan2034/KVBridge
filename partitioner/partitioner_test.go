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
		Address:           ":6379",
		Grpc_address:      "localhost:50051",
		LogPath:           "stdout",
		DataPath:          "./testing/storage",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: 2,
	}
	c2 := &config.Config{
		Address:           ":6380",
		Grpc_address:      "localhost:50052",
		LogPath:           "stdout",
		DataPath:          "./testing/storage2",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: 2,
	}
	c3 := &config.Config{
		Address:           ":6381",
		Grpc_address:      "localhost:50053",
		LogPath:           "stdout",
		DataPath:          "./testing/storage3",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: 2,
	}

	// Create and launch both KVNodes on separate goroutines
	var err error
	node1, cancelFunc1, err := node.NewKVNode(c1)
	if err != nil {
		t.Errorf("node1 creation failed: %v", err)
	}
	node2, cancelFunc2, err := node.NewKVNode(c2)
	if err != nil {
		t.Errorf("node2 creation failed: %v", err)
	}
	node3, cancelFunc3, err := node.NewKVNode(c3)
	if err != nil {
		t.Errorf("node3 creation failed: %v", err)
	}

	// Launch both servers
	go func() {
		err := node1.Start()
		if err != nil {
			t.Errorf("node1 failed: %v", err)
		}
	}()

	go func() {
		err := node2.Start()
		if err != nil {
			t.Errorf("node2 failed: %v", err)
		}
	}()

	go func() {
		err := node3.Start()
		if err != nil {
			t.Errorf("node3 failed: %v", err)
		}
	}()

	// Wait for servers to come up
	time.Sleep(1000 * time.Millisecond)

	return func(t *testing.T) {
		os.RemoveAll("./testing")
		cancelFunc1()
		cancelFunc2()
		cancelFunc3()
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
		t.Errorf("Unknown NodeID assigned: %s", p)
	}
}

func TestConsistentHashPartitioner_GetPartitions(t *testing.T) {
	teardownTest, node1, _, _ := setupTest(t)
	defer teardownTest(t)

	partitions, err := node1.Partitioner.GetPartitions([]byte("testKey"))
	replicationFactor := 2

	if err != nil {
		t.Errorf("Test failed: %s", err)
	}
	if len(partitions) != replicationFactor {
		t.Errorf("Received %d nodes in preference list, expected %d", len(partitions), replicationFactor)
	}

}
