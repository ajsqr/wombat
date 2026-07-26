package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

const backOff = 5

type Tunnel struct {
	config     *config.AgentTunnelConfig
	logger     *slog.Logger
	caCert     []byte
	serverName string
}

func NewTunnel(config *config.AgentTunnelConfig, caCert []byte, serverName string, logger *slog.Logger) *Tunnel {
	return &Tunnel{
		config:     config,
		logger:     logger,
		caCert:     caCert,
		serverName: serverName,
	}
}
func (t *Tunnel) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := t.run()
		if err != nil {
			t.logger.Error("tunnel connection failed", slog.Any("error", err))
		}

		time.Sleep(time.Second * backOff)
	}
}
func (t *Tunnel) run() error {
	t.logger.Info("attempting to connect to the tunnel")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(t.caCert)
	cfg := &tls.Config{
		RootCAs:    pool,
		ServerName: t.serverName,
		MinVersion: tls.VersionTLS13,
	}

	conn, err := tls.Dial("tcp", t.config.Tunnel, cfg)
	if err != nil {
		return fmt.Errorf("attempting to establish a tunnel : %w", err)
	}

	frameWriter := frame.NewWriter(conn)
	frameReader := frame.NewReader(conn)
	t.logger.Info("authenticating tunnel")
	err = t.handshake(frameWriter, t.config)
	if err != nil {
		return fmt.Errorf("error during handshake : %w", err)

	}

	t.logger.Info("successfully authenticated tunnel connection")
	sessionStore := session.NewSessionStore()
	endpoint := NewClientHandler(sessionStore, t.config.Local, t.logger)
	tunnel := tunnel.NewTunnel(conn, frameWriter, frameReader, sessionStore, endpoint)
	t.logger.Info("successfully established tunnel connection")
	err = tunnel.Stream()
	if err != nil {
		return fmt.Errorf("tunnel streaming failed : %w", err)
	}

	return nil
}
