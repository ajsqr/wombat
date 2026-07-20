package config

type AgentConfig struct {
	Tunnels []*AgentTunnelConfig `json:"tunnels"`
}

type AgentTunnelConfig struct {
	// Tunnel is the ip address of the wombat server. A client initiates the tunnel connection to this address.
	// It should be of the format host:port
	// example 192.168.1.15:1234
	Tunnel string `json:"tunnel"`
	// LocalServer is the ip address of the locally running application server.
	// It should be of the format host:port
	// example 127.0.0.1:1234
	Local string `json:"local"`
	// Name is used to easily identify a tunnel
	Name string `json:"name"`
}
