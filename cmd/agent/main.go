package main

import "github.com/ajsqr/wombat/client"

func main() {
	agent := client.Agent{
		TunnelAddr:  "13.234.111.200:8123",
		LocalServer: "127.0.0.1:1232",
	}

	err := agent.Run()
	panic(err)
}
