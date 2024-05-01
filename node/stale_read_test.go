package node_test

import (
	"KVBridge/config"
	"KVBridge/node"
	"KVBridge/types"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
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
	log.Printf("ToxiProxy being initialized")
}

func BenchmarkNode_StaleReadTest(t *testing.B) {

	// start node 1, 2, 3
	os.RemoveAll("./testing")

	// Define configs for both nodes
	c1 := config.DefaultConfig()
	c1.Timeout = 1000 * time.Millisecond
	c1.Grpc_address = "localhost:50051"
	c1.LogPath = ""
	c1.DataPath = "./testing/storage1"
	c1.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}
	c1.WritePreference = types.OpLocal
	c1.ReadPreference = types.OpLocal
	c1.SetWriteThreshold()
	c1.SetReadThreshold()

	// "reference" node -- fetches from all
	c2 := config.DefaultConfig()
	c2.Address = ":6380"
	c2.Grpc_address = "localhost:50052"
	c2.LogPath = ""
	c2.DataPath = "./testing/storage2"
	c2.Timeout = 1000 * time.Millisecond
	c2.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}
	c2.WritePreference = types.OpLocal
	c2.ReadPreference = types.OpLocal
	c2.SetWriteThreshold()
	c2.SetReadThreshold()

	c3 := config.DefaultConfig()
	c3.Address = ":6381"
	c3.Grpc_address = "localhost:50053"
	c3.LogPath = ""
	c3.DataPath = "./testing/storage3"
	c3.Timeout = 1000 * time.Millisecond
	c3.BootstrapServers = []string{"localhost:50151", "localhost:50152", "localhost:50153"}
	c3.WritePreference = types.OpLocal
	c3.ReadPreference = types.OpLocal
	c3.SetWriteThreshold()
	c3.SetReadThreshold()

	// Create and launch all KVNodes on separate goroutines
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
	time.Sleep(2000 * time.Millisecond)

	teardownTest := func(t *testing.B) {
		//os.RemoveAll("./testing")
		cancelFunc1()
		cancelFunc2()
		cancelFunc3()
	}

	defer teardownTest(t)

	// set up toxiproxy

	initToxiProxy()
	// Okay actual test starts here lol

	// Setup "client" reads and writes
	//
	// Define the range of keys that will be picked
	maxKey := int64(100000)

	// note: values will always be between those too as well
	//latestValue := make(map[string]string)

	// communication channel for keys, errors
	errChan := make(chan error, 1000)
	resChan := make(chan int, 1000)

	writeFunc := func(key string, value string) {
		// Write this to node 3 always
		errChan <- node3.Write([]byte(key), []byte(value))
	}

	readFunc := func(key, value string) {
		counter := 0
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		timeout := time.NewTimer(5 * time.Second)
		defer timeout.Stop()
		for {
			select {
			case <-ticker.C:
				// Read this from node 1 always
				res, err := node1.Read([]byte(key))
				if err != nil {
					errChan <- err
				}
				if string(res) == (value) {
					resChan <- counter
					return
				}
				counter++

			case <-timeout.C:
				resChan <- -1
				return
			}
		}
	}

	// mutex for the map
	//var mu sync.Mutex

	// looper
	looper := func() {
		var wg sync.WaitGroup
		// generate a random key between the given ranges
		keyInt := rand.Int64N(maxKey)
		valueInt := rand.Int64()
		key := strconv.FormatInt(keyInt, 10)
		value := strconv.FormatInt(valueInt, 10)
		//mu.Lock()
		//latestValue[key] = value
		//mu.Unlock()

		wg.Add(1)
		go func(key, value string) {
			defer wg.Done()
			writeFunc(key, value)
		}(key, value)

		wg.Add(1)
		go func(key, value string) {
			defer wg.Done()
			readFunc(key, value)
		}(key, value)

		wg.Wait()
	}

	// call looper in a loop async
	maxParallel := 1
	maxRequests := 1000

	parallelChan := make(chan struct{}, maxParallel)

	var wg sync.WaitGroup
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < maxRequests; i++ {
		//t.Logf("request: %d", i)
		parallelChan <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			looper()
		}()
		<-parallelChan
		//t.Logf("done: %d", i)
	}

	// stat counter
	numStaleReads := 0
	numRequests := 0
	numErrs := 0
	var countWg sync.WaitGroup
	countWg.Add(1)
	go func() {
		defer countWg.Done()
		for c := range resChan {
			numRequests++
			if c == -1 {
				// this was a timeout
				numErrs++
			} else {
				numStaleReads += c
			}
		}
	}()

	// error printer
	countWg.Add(1)
	go func() {
		defer countWg.Done()
		for err := range errChan {
			// this was a write error
			if err != nil {
				t.Logf("error: %v", err)
			}
		}
	}()

	// wait for looper
	wg.Wait()
	close(resChan)
	close(errChan)

	// wait for counters
	countWg.Wait()

	t.Logf("Avg Staleness: %f\tnumRequests: %d\tnumStaleReads: %d\tnumErrs: %d\t", float32(numStaleReads)/float32(numRequests), numRequests, numStaleReads, numErrs)
}
