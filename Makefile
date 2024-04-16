run:
	go run main.go

build: 
	go build -o kvbridge main.go

clean:
	rm -rf ./kvbridge
