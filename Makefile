build-server:
	go build -o build/server ./cmd/server
build-agent:
	go build -o build/agent ./cmd/agent
build: build-server build-agent

build-server-amd:
	GOOS=linux GOARCH=amd64 go build -o build/server ./cmd/server
build-agent-amd:
	GOOS=linux GOARCH=amd64 go build -o build/agent ./cmd/agent
build-amd: build-server-amd build-agent-amd 