package server

import (
	"log/slog"
	"sync"
)

type ServerConfig struct {
	Channels []*ChannelConfig `json:"channels"`
}

type Server struct {
	config *ServerConfig
	logger *slog.Logger
}

func NewServer(config *ServerConfig, logger *slog.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
	}
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	for _, channelConfig := range s.config.Channels {
		logger := s.logger.With(slog.String("channel", channelConfig.Name))
		wg.Add(1)
		channel := NewChannel(channelConfig, logger)
		go channel.Run(&wg)
	}

	wg.Wait()
}
