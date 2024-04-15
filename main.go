package main

// Entry point for a backend node, which starts relevant services.
// For now it opens up a simple backend server that listens for a client.
func main() {
	node := DefaultConfig().Build()
	node.Start()
}
