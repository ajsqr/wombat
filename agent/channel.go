package agent

import (
	"log/slog"
	"net"
	"sync"

	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/receiver/session"
)

type ChannelConfig struct {
	// TunnelAddr is the ip address of the wombat server. A client initiates the tunnel connection to this address.
	// It should be of the format host:port
	// example 192.168.1.15:1234
	TunnelAddr string `json:"tunnelAddr"`
	// LocalServer is the ip address of the locally running application server.
	// It should be of the format host:port
	// example 127.0.0.1:1234
	LocalServer string `json:"localServer"`
	// Name is used to easily identify a tunnel
	Name string `json:"name"`
}

type Channel struct {
	config *ChannelConfig
	logger *slog.Logger
}

func NewChannel(config *ChannelConfig, logger *slog.Logger) *Channel {
	return &Channel{
		config: config,
		logger: logger,
	}
}

func (c *Channel) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	c.logger.Info("attempting to connect to the tunnel")
	conn, err := net.Dial("tcp", c.config.TunnelAddr)
	if err != nil {
		c.logger.Error("attempting to establish a tunnel", slog.Any("error", err))
		return
	}

	c.logger.Info("successfully connected to the tunnel")

	sessionStore := session.NewSessionStore()
	endpoint := NewClientHandler(sessionStore, c.config.LocalServer, c.logger)
	tunnel := tunnel.NewTunnel(conn, sessionStore, endpoint)
	err = tunnel.Stream()
	if err != nil {
		c.logger.Error("tunnel streaming failed", slog.Any("error", err))
		return
	}
}
