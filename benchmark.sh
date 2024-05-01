#!/bin/bash

# set up proxies
toxiproxy-cli create -l localhost:50151 -u localhost:50051 node1
toxiproxy-cli create -l localhost:50152 -u localhost:50052 node2
toxiproxy-cli create -l localhost:50153 -u localhost:50053 node3

# add toxics
toxiproxy-cli toxic add --type latency --upstream -a latency=10 -a jitter=2 node1
toxiproxy-cli toxic add --type latency --upstream -a latency=10 -a jitter=2 node2
toxiproxy-cli toxic add --type latency --upstream -a latency=10 -a jitter=2 node3

# build
make build

# start kvbridge nodes
./kvbridge --conf config.yaml &
./kvbridge --conf config2.yaml &
./kvbridge --conf config3.yaml &

sleep 1

# run benchmark
redis-benchmark -t set,get -n 1000 -q -c 4 # --csv

# terminate nodes
redis-cli -p 6379 -n 0 close
redis-cli -p 6380 -n 0 close
redis-cli -p 6381 -n 0 close

# delete proxies
toxiproxy-cli delete node1
toxiproxy-cli delete node2
toxiproxy-cli delete node3
