package server

import (
	"log/slog"
	"sync"

	"github.com/ajsqr/wombat/config"
)

type Server struct {
	config *config.ServerConfig
	logger *slog.Logger
}

func NewServer(logger *slog.Logger) (*Server, error) {
	var serverConfig config.ServerConfig
	err := config.LoadServerConfig(&serverConfig)
	if err != nil {
		return nil, err
	}
	return &Server{
		config: &serverConfig,
		logger: logger,
	}, nil
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	for _, tunnelConfig := range s.config.Tunnels {
		logger := s.logger.With(slog.String("channel", tunnelConfig.Name))
		wg.Add(1)
		channel := NewTunnel(tunnelConfig, logger)
		go channel.Run(&wg)
	}

	wg.Wait()
}
