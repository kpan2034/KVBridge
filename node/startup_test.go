package node_test

import (
	"KVBridge/config"
	"KVBridge/node"
	"KVBridge/proto/compiled/startup"
	. "KVBridge/types"
	"context"
	"fmt"
	"math"
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

func setupTest2(t *testing.T, rf int) (teardownTest func(t *testing.T), node1, node2, node3 *node.KVNode, f1, f2, f3 func()) {
	os.RemoveAll("./testing")

	// Define configs for both nodes
	c1 := &config.Config{
		Address:           ":6379",
		Grpc_address:      "localhost:50051",
		LogPath:           "stdout",
		DataPath:          "./testing/storage",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: rf,
	}
	c2 := &config.Config{
		Address:           ":6380",
		Grpc_address:      "localhost:50052",
		LogPath:           "stdout",
		DataPath:          "./testing/storage2",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: rf,
	}
	c3 := &config.Config{
		Address:           ":6381",
		Grpc_address:      "localhost:50053",
		LogPath:           "stdout",
		DataPath:          "./testing/storage3",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: rf,
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
	}, node1, node2, node3, cancelFunc1, cancelFunc2, cancelFunc3
}

func TestMessager_StartupService(t *testing.T) {
	teardownTest, node1, _, _ := setupTest(t)
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

func TestMessager_Recover(t *testing.T) {
	_, node1, _, _, cancelFunc1, cancelFunc2, cancelFunc3 := setupTest2(t, 2)

	err := node1.Write([]byte("testKey1"), []byte("testVal1"))
	if err != nil {
		t.Errorf("TestMessager_Recover failed %s", err)
	}

	iter, err := node1.Storage.GetSnapshotIter(0, math.MaxUint32)
	if err != nil {
		t.Errorf("TestMessager_Recover failed %s", err)
	}
	iter.First()
	if !iter.Valid() {
		fmt.Printf("iterator not valid")
	} else {
		fmt.Printf("^^^^ %v ^^^^", iter.Key())
	}

	//cancelFunc1()
	//err = node2.Write([]byte("testKey2"), []byte("testVal2"))
	//if err != nil {
	//	t.Errorf("TestMessager_Recover failed %s", err)
	//}
	//
	//c1 := &config.Config{
	//	Address:           ":6379",
	//	Grpc_address:      "localhost:50051",
	//	LogPath:           "stdout",
	//	DataPath:          "./testing/storage",
	//	BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
	//	ReplicationFactor: 3,
	//}
	////node1_rec, node1RecCancelFunc, err := node.NewKVNode(c1)
	//node1_rec, node1RecCancelFunc, err := node.NewKVNode(c1)
	//
	//if err != nil {
	//	t.Errorf("node1 creation failed: %v", err)
	//}
	//go func() {
	//	err := node1_rec.Start()
	//	if err != nil {
	//		t.Errorf("node1 failed: %v", err)
	//	}
	//}()
	//
	//// Wait for servers to come up
	//time.Sleep(1000 * time.Millisecond)
	//
	//err = node1_rec.Recover()
	//
	//if err != nil {
	//	t.Errorf("Recover failed: %s", err)
	//}
	//
	//err = os.RemoveAll("./testing")
	//if err != nil {
	//	t.Errorf("Recover failed: %s", err)
	//}
	//node1RecCancelFunc()
	cancelFunc2()
	cancelFunc3()
	cancelFunc1()
}
