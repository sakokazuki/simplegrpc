build:
	GO111MODULE=on go build -tags=release
dev:
	go run main.go
client: 
	go run client/main.go