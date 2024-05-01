package node

import (
	"KVBridge/config"
	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"log"
	"os"
	"testing"
	"time"
)

var toxiClient *toxiproxy.Client
var proxies []*toxiproxy.Proxy

func initToxiProxy() {
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
	latency := 10            // milliseconds
	jitter := 2              // milliseconds

	for _, p := range proxies {
		_, err := p.AddToxic("latency_down", "latency", "upstream", toxicity, toxiproxy.Attributes{
			"latency": latency,
			"jitter":  jitter,
		})
		if err != nil {
			log.Printf("Toxiproxy failed to add toxic %s", err)
		}
		err = p.Enable()
		if err != nil {
			return
		}
	}
}

func SetupTestEmulation(t *testing.T) (teardownTest func(t *testing.T), node1, node2, node3 *KVNode, f1, f2, f3 func()) {
	os.RemoveAll("./testing")
	initToxiProxy()

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
	c1.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}

	c3 := config.DefaultConfig()
	c3.Address = ":6381"
	c3.Grpc_address = "localhost:50053"
	c3.LogPath = "stdout"
	c3.DataPath = "./testing/storage3"
	c3.Timeout = 10000 * time.Millisecond
	c1.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}

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
