package config

type ServerConfig struct {
	Tunnels []*ServerTunnelConfig `json:"tunnels"`
}

type ServerTunnelConfig struct {
	// Tunnel is the ip address of the wombat server. A client initiates the tunnel connection to this address.
	// It should be of the format host:port
	// example 192.168.1.15:1234
	Tunnel    string `json:"tunnel"`
	Public    string `json:"public"`
	Name      string `json:"name"`
	TokenName string `json:"tokenName"`
}
