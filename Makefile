run:
	go run main.go node.go interfaces.go config.go

build: 
	go build -o kvbridge main.go node.go interfaces.go config.go 

