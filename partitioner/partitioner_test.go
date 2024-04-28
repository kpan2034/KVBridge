package partitioner_test

import (
	"KVBridge/config"
	"KVBridge/node"
	"os"
	"strconv"
	"testing"
	"time"
)

func setupTest(numNodes int, replicationFactor int, t *testing.T, testFolder string, addrStartPort int, grpcStartPort int) (func(t *testing.T), []*node.KVNode) {
	os.RemoveAll(testFolder)

	nodes := make([]*node.KVNode, numNodes)

	bootstrapServers := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		bootstrapServers[i] = "localhost:" + strconv.FormatInt(int64(grpcStartPort+i), 10)
	}

	cancelFuncs := make([]func(), numNodes)

	for i := 0; i < numNodes; i++ {
		c := &config.Config{
			Address:           ":" + strconv.FormatInt(int64(addrStartPort+i), 10),
			Grpc_address:      "localhost:" + strconv.FormatInt(int64(grpcStartPort+i), 10),
			LogPath:           "stdout",
			DataPath:          testFolder + "/storage" + strconv.FormatInt(int64(i), 10),
			BootstrapServers:  bootstrapServers,
			ReplicationFactor: replicationFactor,
		}
		newNode, cancelFunc, err := node.NewKVNode(c)
		if err != nil {
			t.Errorf("node%d creation failed: %v", i, err)
		}
		go func() {
			err := newNode.Start()
			if err != nil {
				t.Errorf("node%d startup failed: %v", i, err)
			}
		}()
		nodes[i] = newNode
		cancelFuncs[i] = cancelFunc
	}

	// Wait for servers to come up
	time.Sleep(1000 * time.Millisecond)

	return func(t *testing.T) {
		os.RemoveAll(testFolder)
		for _, f := range cancelFuncs {
			f()
		}
	}, nodes
}

func TestSimplePartitioner_GetPartitions(t *testing.T) {
	teardownTest, nodes := setupTest(3, 1, t, "./testing/test1", 6379, 50051)
	node1 := nodes[0]
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
	replicationFactor := 3
	teardownTest, nodes := setupTest(9, replicationFactor, t, "./testing/test2", 6390, 50060)
	defer teardownTest(t)

	partitions, err := nodes[0].Partitioner.GetPartitions([]byte("testKey"))

	if err != nil {
		t.Errorf("Test failed: %s", err)
	}
	if len(partitions) != replicationFactor {
		t.Errorf("Received %d nodes in preference list, expected %d", len(partitions), replicationFactor)
	}

}
