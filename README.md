# KVBridge, a Practical Key-Value store.

## Basic roadmap:

1. Have Key-Value APIs up and running for a single node. Storage can be non-persistent, in memory.

Should we design APIs using gRPC or RESP?

Basic commands to support (taken from Redis)

- SET
- GET
- DEL
- PING

- QUIT
- PUB
- SUB

We could go with [this library](https://github.com/tidwall/redcon) for client-server communication.

2. Set up storage engine. Probably using [RocksDB](https://github.com/facebook/rocksdb.git) or [FASTER]()
3. Set up communication between nodes. Only basic client redirection handling. Single leader, no fail-over, no consensus, no replication.

4. Add basic replication between nodes, using [this Raft library](https://github.com/etcd-io/raft):

5. Add read-replica support.

6.
