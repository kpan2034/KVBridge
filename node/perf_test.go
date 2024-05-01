package node_test

import (
	"KVBridge/config"
	. "KVBridge/node"
	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var latencyNode1 int = 500 // milliseconds
var jitterNode1 int = 2    // milliseconds
var latency int = 10       // milliseconds
var jitter int = 2         // milliseconds
var num_gets = 100
var num_sets = 100

func initToxiProxyPerf() {
	if toxiClient != nil {
		// since we only initialize once
		return
	}
	log.Printf("ToxiProxy being initialized")
	toxiClient = toxiproxy.NewClient("localhost:8474")

	p1, err := toxiClient.CreateProxy("node1", "localhost:50151", "localhost:50051")
	if err != nil {
		return
	}
	err = p1.Enable()
	if err != nil {
		return
	}
	p2, err := toxiClient.CreateProxy("node2", "localhost:50152", "localhost:50052")
	if err != nil {
		return
	}
	err = p2.Enable()
	if err != nil {
		return
	}
	p3, err := toxiClient.CreateProxy("node3", "localhost:50153", "localhost:50053")
	if err != nil {
		return
	}
	err = p3.Enable()
	if err != nil {
		return
	}
	proxies = []*toxiproxy.Proxy{p1, p2, p3}
	toxicity := float32(1.0) // Ratio of requests to which proxy rules are applied

	_, err = p1.AddToxic("latency_down", "latency", "upstream", toxicity, toxiproxy.Attributes{
		"latency": latencyNode1,
		"jitter":  jitterNode1,
	})
	if err != nil {
		log.Printf("Toxiproxy failed to add toxic %s", err)
	}
	err = p1.Enable()
	if err != nil {
		return
	}

	_, err = p2.AddToxic("latency_down", "latency", "upstream", toxicity, toxiproxy.Attributes{
		"latency": latency,
		"jitter":  jitter,
	})
	if err != nil {
		log.Printf("Toxiproxy failed to add toxic %s", err)
	}
	err = p2.Enable()
	if err != nil {
		return
	}

	_, err = p3.AddToxic("latency_down", "latency", "upstream", toxicity, toxiproxy.Attributes{
		"latency": latency,
		"jitter":  jitter,
	})
	if err != nil {
		log.Printf("Toxiproxy failed to add toxic %s", err)
	}
	err = p3.Enable()
	if err != nil {
		return
	}
}

func SetupTestEmulation2(t *testing.T) (teardownTest func(t *testing.T), node1, node2, node3 *KVNode, f1, f2, f3 func()) {
	os.RemoveAll("./testing")
	initToxiProxyPerf()

	// Define configs for both nodes
	c1 := config.DefaultConfig()
	c1.Timeout = 10000 * time.Millisecond
	c1.Grpc_address = "localhost:50051"
	c1.LogPath = "stdout"
	c1.DataPath = "./testing/storage1"
	c1.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}

	c2 := config.DefaultConfig()
	c2.Address = ":6380"
	c2.Grpc_address = "localhost:50052"
	c2.LogPath = "stdout"
	c2.DataPath = "./testing/storage2"
	c2.Timeout = 10000 * time.Millisecond
	c2.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}

	c3 := config.DefaultConfig()
	c3.Address = ":6381"
	c3.Grpc_address = "localhost:50053"
	c3.LogPath = "stdout"
	c3.DataPath = "./testing/storage3"
	c3.Timeout = 10000 * time.Millisecond
	c3.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}

	// Create and launch both KVNodes on separate goroutines
	var err error
	node1, cancelFunc1, err := NewKVNode(c1)
	if err != nil {
		t.Errorf("node1 creation failed: %v", err)
	}
	node2, cancelFunc2, err := NewKVNode(c2)
	if err != nil {
		t.Errorf("node2 creation failed: %v", err)
	}
	node3, cancelFunc3, err := NewKVNode(c3)
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

func TestSlowNode1(t *testing.T) {
	_, node1, node2, node3, cancelFunc1, cancelFunc2, cancelFunc3 := SetupTestEmulation2(t)
	rand.Seed(0)

	time.Sleep(1000)
	write_failures := 0
	node1TimeElapsed := int64(0)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		startTime := time.Now()
		for i := 0; i < num_sets; i++ {
			k := "testKey" + strconv.FormatInt(int64(i/2), 10)
			err := node1.Write([]byte(k), []byte("testVal1"))
			if err != nil {
				write_failures += 1
			}
		}
		timeElapsedMS := (time.Now().Sub(startTime)).Milliseconds()
		node1TimeElapsed = timeElapsedMS
		//t.Logf("STATS# Node 1: %d write reqs in %v ms | %v rps", num_sets, timeElapsedMS, float64(num_sets)/(float64(timeElapsedMS)/1000))
	}()

	read_failures_node2 := 0
	node2TimeElapsed := int64(0)
	go func() {
		defer wg.Done()
		startTime := time.Now()
		for i := 0; i < (num_gets / 2); i++ {
			k := "testKey" + strconv.FormatInt(int64(i), 10)
			_, err := node2.Read([]byte(k))
			if err != nil {
				read_failures_node2 += 1
			}
		}
		timeElapsedMS := (time.Now().Sub(startTime)).Milliseconds()
		node2TimeElapsed = timeElapsedMS
		//t.Logf("STATS# Node 2: %d read reqs in %v ms | %v rps", num_gets/2, timeElapsedMS, float64(num_gets/2)/(float64(timeElapsedMS)/1000))
	}()
	read_failures_node3 := 0
	node3TimeElapsed := int64(0)
	go func() {
		defer wg.Done()
		startTime := time.Now()
		for i := 0; i < (num_gets / 2); i++ {
			k := "testKey" + strconv.FormatInt(int64(i), 10)
			_, err := node3.Read([]byte(k))
			if err != nil {
				read_failures_node3 += 1
			}
		}
		timeElapsedMS := (time.Now().Sub(startTime)).Milliseconds()
		node3TimeElapsed = timeElapsedMS
		//t.Logf("STATS# Node 3: %d read reqs in %v ms | %v rps", num_gets/2, timeElapsedMS, float64(num_gets/2)/(float64(timeElapsedMS)/1000))
	}()
	wg.Wait()
	//time.Sleep(60 * 2 * time.Second)
	t.Logf("STATS# Node1 %d/%d writes failed\n Node2 %d/%d reads failed\n Node3 %d/%d reads failed",
		write_failures, num_sets, read_failures_node2, num_gets/2, read_failures_node3, num_gets/2)
	t.Logf("STATS# Node 1(latency: %dms): %d write reqs in %v ms | %v rps", latencyNode1, num_sets, node1TimeElapsed, float64(num_sets)/(float64(node1TimeElapsed)/1000))
	t.Logf("STATS# Node 2(latency: %dms): %d read reqs in %v ms | %v rps", latency, num_gets/2, node2TimeElapsed, float64(num_gets/2)/(float64(node2TimeElapsed)/1000))
	t.Logf("STATS# Node 3(latency: %dms): %d read reqs in %v ms | %v rps", latency, num_gets/2, node3TimeElapsed, float64(num_gets/2)/(float64(node3TimeElapsed)/1000))

	cancelFunc1()
	cancelFunc2()
	cancelFunc3()
}
