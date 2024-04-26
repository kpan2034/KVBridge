package main

import (
	"KVBridge/config"
	"KVBridge/node"
	"fmt"
	"os"
)

// Entry point for a backend node, which starts relevant services.
// For now it opens up a simple backend server that listens for a client.
func main() {
	config := config.NewConfigFromEnv()
	kvNode, err := node.NewKVNode(config)
	if err != nil {
		fmt.Printf("Failed to create KV node: %s", err)
		os.Exit(1)
	}
	err = kvNode.Start()
	if err != nil {
		fmt.Printf("Failed to start KV node: %s", err)
		os.Exit(1)
	}
}
