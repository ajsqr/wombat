package main

import "github.com/ajsqr/wombat/server"

func main() {
	cfg := server.Config{
		TunnelAddr: "0.0.0.0:8123",
		ServerAddr: "0.0.0.0:8000",
	}

	server := server.NewServer(&cfg)
	err := server.Run()
	panic(err)
}
