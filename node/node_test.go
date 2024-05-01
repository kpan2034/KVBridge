package node_test

import (
	"KVBridge/config"
	"KVBridge/node"
	"context"
	"math"
	"os"
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
	teardownTest, node1, _, _, _, _, _ := setupTest2(t, 3)
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

// Recovery tests
func setupTest2(t *testing.T, rf int) (teardownTest func(t *testing.T), node1, node2, node3 *node.KVNode, f1, f2, f3 func()) {
	os.RemoveAll("./testing")

	// Define configs for both nodes
	c1 := config.DefaultConfig()
	c1.Timeout = 10000 * time.Millisecond
	c1.LogPath = "stdout"
	c1.DataPath = "./testing/storage1"

	c2 := config.DefaultConfig()
	c2.Address = ":6380"
	c2.Grpc_address = "localhost:50052"
	c2.LogPath = "stdout"
	c2.DataPath = "./testing/storage2"
	c2.Timeout = 10000 * time.Millisecond

	c3 := config.DefaultConfig()
	c3.Address = ":6381"
	c3.Grpc_address = "localhost:50053"
	c3.LogPath = "stdout"
	c3.DataPath = "./testing/storage3"
	c3.Timeout = 10000 * time.Millisecond

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
	}, node1, node2, node3, cancelFunc1, cancelFunc2, cancelFunc3
}

func TestMessager_Recover(t *testing.T) {
	_, node1, node2, _, cancelFunc1, cancelFunc2, cancelFunc3 := setupTest2(t, 3)

	err := node1.Write([]byte("testKey1"), []byte("testVal1"))
	if err != nil {
		t.Errorf("TestMessager_Recover failed %s", err)
	}
	if getNumRecords(node1, t) != 1 || getNumRecords(node2, t) != 1 {
		t.Errorf("Node1 Records: Expected %d Actual %d \nNode2 Records: Expected %d Actual %d",
			1, getNumRecords(node1, t), 1, getNumRecords(node2, t))
	}

	cancelFunc1()
	err = node2.Write([]byte("testKey2"), []byte("testVal2"))
	if err != nil {
		// TODO: node2 currently tries to connect to node 1 and fails!
		t.Logf("TestMessager_Recover failed %s", err)
	}

	c1 := config.DefaultConfig()
	c1.Timeout = 10000 * time.Millisecond
	c1.LogPath = "stdout"
	c1.DataPath = "./testing/storage1"
	node1_rec, node1RecCancelFunc, err := node.NewKVNode(c1)

	go func() {
		err := node1_rec.Start()
		if err != nil {
			t.Errorf("node1 failed: %v", err)
		}
	}()

	// Wait for servers to come up
	time.Sleep(1000 * time.Millisecond)

	if getNumRecords(node1_rec, t) != 1 || getNumRecords(node2, t) != 2 {
		t.Errorf("Node1 (Before Recover) Records: Expected %d Actual %d \nNode2 Records: Expected %d Actual %d",
			1, getNumRecords(node1_rec, t), 2, getNumRecords(node2, t))
	}

	err = node1_rec.Recover()
	if err != nil {
		t.Errorf("Recover failed: %s", err)
	}
	if getNumRecords(node1_rec, t) != 2 || getNumRecords(node2, t) != 2 {
		t.Errorf("Node1(Recovered) Records: Expected %d Actual %d \nNode2 Records: Expected %d Actual %d",
			2, getNumRecords(node1_rec, t), 2, getNumRecords(node2, t))
	}

	err = os.RemoveAll("./testing")
	if err != nil {
		t.Errorf("Recover failed: %s", err)
	}
	node1RecCancelFunc()
	cancelFunc2()
	cancelFunc3()
}

func getNumRecords(node *node.KVNode, t *testing.T) int {
	snapshotDB := node.Storage.GetSnapshotDB()
	iters, err := node.Storage.GetSnapshotIters(0, math.MaxUint32, snapshotDB)
	if err != nil {
		t.Errorf("getNumRecord failed %s", err)
	}

	count := 0
	for _, iter := range iters {
		iter.First()
		for iter.Valid() {
			count += 1
			iter.Next()
		}
		err := iter.Close()
		if err != nil {
			t.Errorf("getNumRecord failed %s", err)
		}
	}
	err = snapshotDB.Close()
	if err != nil {
		t.Errorf("getNumRecord failed %s", err)
	}
	return count
}
