package server

import (
	"crypto/tls"
	"log/slog"
	"sync"

	"github.com/ajsqr/wombat/config"
)

type Server struct {
	cert   tls.Certificate
	config *config.ServerConfig
	logger *slog.Logger
}

func NewServer(logger *slog.Logger) (*Server, error) {
	var serverConfig config.ServerConfig
	err := config.LoadServerConfig(&serverConfig)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(serverConfig.CertPath, serverConfig.KeyPath)
	if err != nil {
		return nil, err
	}

	return &Server{
		config: &serverConfig,
		logger: logger,
		cert:   cert,
	}, nil
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	for _, tunnelConfig := range s.config.Tunnels {
		logger := s.logger.With(slog.String("channel", tunnelConfig.Name))
		wg.Add(1)
		channel := NewTunnel(tunnelConfig, s.cert, logger)
		go channel.Run(&wg)
	}

	wg.Wait()
}
