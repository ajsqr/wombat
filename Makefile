build-server:
	go build -o build/server ./cmd/server/main.go 
build-agent:
	go build -o build/agent ./cmd/agent/main.go
build: build-server build-agent

build-server-amd:
	GOOS=linux GOARCH=amd64 go build -o build/server ./cmd/server/main.go 
build-agent-amd:
	GOOS=linux GOARCH=amd64 go build -o build/agent ./cmd/agent/main.go
build-amd: build-server-amd build-agent-amd 