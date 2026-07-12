package main

import "github.com/ajsqr/wombat/server"

func main() {
	cfg := server.Config{
		TunnelAddr: "127.0.0.1:1231",
		ServerAddr: "127.0.0.1:1232",
	}

	server := server.NewServer(&cfg)
	err := server.Run()
	panic(err)
}
