package main

import "github.com/ajsqr/wombat/client"

func main() {
	agent := client.Agent{
		TunnelAddr:  "127.0.0.1:1231",
		LocalServer: "127.0.0.1:1233",
	}

	err := agent.Run()
	panic(err)
}
