package client

import (
	"log/slog"
	"net"
	"os"

	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/receiver/session"
)

type Agent struct {
	// TunnelAddr is the ip address of the wombat server. A client initiates the tunnel connection to this address.
	// It should be of the format host:port
	// example 192.168.1.15:1234
	TunnelAddr string `json:"tunnelAddr"`
	// LocalServer is the ip address of the locally running application server.
	// It should be of the format host:port
	// example 127.0.0.1:1234
	LocalServer string `json:"localServer"`
}

func (a *Agent) Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("attempting to connect to the tunnel")
	conn, err := net.Dial("tcp", a.TunnelAddr)
	if err != nil {
		return err
	}

	logger.Info("successfully connected to the tunnel")

	sessionStore := session.NewSessionStore()
	endpoint := NewClientHandler(sessionStore, a.LocalServer, logger)
	tunnel := tunnel.NewTunnel(conn, sessionStore, endpoint)
	return tunnel.Stream()
}
