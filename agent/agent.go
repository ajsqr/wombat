package agent

import (
	"log/slog"
	"sync"
)

type AgentConfig struct {
	Channels []*ChannelConfig `json:"channels"`
}

type Agent struct {
	config *AgentConfig
	logger *slog.Logger
}

func NewAgent(config *AgentConfig, logger *slog.Logger) *Agent {
	return &Agent{
		config: config,
		logger: logger,
	}
}

func (a *Agent) Run() {
	wg := sync.WaitGroup{}
	for _, channelConfig := range a.config.Channels {
		logger := a.logger.With(slog.String("channel", channelConfig.Name))
		wg.Add(1)
		channel := NewChannel(channelConfig, logger)
		go channel.Run(&wg)
	}

	wg.Wait()
}
