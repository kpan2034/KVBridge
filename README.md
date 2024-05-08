# KVBridge, a Practical Key-Value store.


## Building
Clone the repository to your system, cd into the directory and run `make`:


```
git clone https://github.com/kpan2034/KVBridge.git
cd KVBridge
make build
```

This creates a binary called `kvbridge` in the same directory.

## Running

Executing the above binary starts a new KVbridge process with the default configuration.
You can specify the configuration file using the `--conf` parameter:

```
./kvbridge --conf config.yaml
```

Use `-h` to see additional parameters that you can specify from the command line:

```
./kvbridge -h
      --address string             address for server to listen for client (default ":6379")
      --bootstrap_servers string   bootstrap servers in the cluster (default "localhost:50051,localhost:50052")
      --conf strings               path to one or more .yaml config files to load (default [./config.yaml])
      --data_path string           path to persistent storage (default "./tmp/storage")
      --grpc_address string        address for server to listen for other nodes (default "localhost:50051")
      --log_path string            path to logfile (default "./tmp/log")
      --read_pref string           read preference for cluster (default "majority")
      --replication_factor int     number of nodes to replicate each data point (default 3)
      --timeout int                default timeout between nodes in ns (default 10000000000)
      --write_pref string          write preference for cluster (default "majority")
```

