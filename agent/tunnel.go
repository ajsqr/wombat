package agent

import (
	"log/slog"
	"net"
	"sync"

	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/receiver/session"
)

type Tunnel struct {
	config *config.AgentTunnelConfig
	logger *slog.Logger
}

func NewTunnel(config *config.AgentTunnelConfig, logger *slog.Logger) *Tunnel {
	return &Tunnel{
		config: config,
		logger: logger,
	}
}

func (t *Tunnel) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	t.logger.Info("attempting to connect to the tunnel")
	conn, err := net.Dial("tcp", t.config.Tunnel)
	if err != nil {
		t.logger.Error("attempting to establish a tunnel", slog.Any("error", err))
		return
	}

	t.logger.Info("successfully connected to the tunnel")

	sessionStore := session.NewSessionStore()
	endpoint := NewClientHandler(sessionStore, t.config.Local, t.logger)
	tunnel := tunnel.NewTunnel(conn, sessionStore, endpoint)
	err = tunnel.Stream()
	if err != nil {
		t.logger.Error("tunnel streaming failed", slog.Any("error", err))
		return
	}
}
