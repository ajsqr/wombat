build-server:
	go build -o build/server ./cmd/server/main.go 
build-agent:
	go build -o build/agent ./cmd/agent/main.go
build-echo:
	go build -o build/echo ./cmd/echo/main.go 
build: build-server build-agent build-echo

build-server-amd:
	GOOS=linux GOARCH=amd64 go build -o build/server ./cmd/server/main.go 
build-agent-amd:
	GOOS=linux GOARCH=amd64 go build -o build/agent ./cmd/agent/main.go
build-echo-amd:
	GOOS=linux GOARCH=amd64 go build -o build/echo ./cmd/echo/main.go 
build-amd: build-server-amd build-agent-amd build-echo-amd