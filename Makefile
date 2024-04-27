REL_PROTO_PATH=./"proto/compiled"
PROTOFILES=./proto/ping.proto ./proto/startup.proto ./proto/replication.proto

run:
	go run main.go

build: 
	go build -o kvbridge main.go

clean:
	rm -rf ./kvbridge
	rm -rf $(REL_PROTO_PATH)

protoc:
	rm -rf $(REL_PROTO_PATH)
	mkdir -p $(REL_PROTO_PATH)
	protoc --go_out=$(REL_PROTO_PATH) --go-grpc_out=$(REL_PROTO_PATH) $(PROTOFILES)
