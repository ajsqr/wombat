package agent

import (
	"log/slog"
	"net"
	"sync"

	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/frame"
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

	frameWriter := frame.NewWriter(conn)
	frameReader := frame.NewReader(conn)
	t.logger.Info("authenticating tunnel")
	err = t.handshake(frameWriter, t.config)
	if err != nil {
		t.logger.Error("error during handshake", slog.Any("error", err))
		return
	}

	t.logger.Info("successfully authenticated tunnel connection")
	sessionStore := session.NewSessionStore()
	endpoint := NewClientHandler(sessionStore, t.config.Local, t.logger)
	tunnel := tunnel.NewTunnel(conn, frameWriter, frameReader, sessionStore, endpoint)
	t.logger.Info("successfully established tunnel connection")
	err = tunnel.Stream()
	if err != nil {
		t.logger.Error("tunnel streaming failed", slog.Any("error", err))
		return
	}

}
