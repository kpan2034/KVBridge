package node

func (node *KVNode) Write(key []byte, value []byte) error {
	err := node.Storage.Set(key, value)
	if err != nil {
		return err
	}
	node.Logger.Debugf("wrote (%v:%v) to local storage", key, value)

	// replicate write to other node
	nacks, err := node.ReplicateWrites(key, value)
	if err != nil {
		// just error out for now, later should just log
		return err
	}
	node.Logger.Debugf("replicated to %d other nodes", nacks)

	return nil
}
