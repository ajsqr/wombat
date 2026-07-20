package agent

import (
	"log/slog"
	"sync"

	"github.com/ajsqr/wombat/config"
)

type Agent struct {
	config *config.AgentConfig
	logger *slog.Logger
}

func NewAgent(logger *slog.Logger) (*Agent, error) {
	var agentConfig config.AgentConfig
	err := config.LoadAgentConfig(&agentConfig)
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: &agentConfig,
		logger: logger,
	}, nil
}

func (a *Agent) Run() {
	wg := sync.WaitGroup{}
	for _, channelConfig := range a.config.Tunnels {
		logger := a.logger.With(slog.String("channel", channelConfig.Name))
		wg.Add(1)
		channel := NewTunnel(channelConfig, logger)
		go channel.Run(&wg)
	}

	wg.Wait()
}
