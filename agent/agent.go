package agent

import (
	"log/slog"
	"os"
	"sync"

	"github.com/ajsqr/wombat/config"
)

type Agent struct {
	config     *config.AgentConfig
	logger     *slog.Logger
	serverName string
	caCert     []byte
}

func NewAgent(logger *slog.Logger) (*Agent, error) {
	var agentConfig config.AgentConfig
	err := config.LoadAgentConfig(&agentConfig)
	if err != nil {
		return nil, err
	}

	ca, err := os.ReadFile(agentConfig.CACertPath)
	if err != nil {
		return nil, err
	}

	return &Agent{
		config:     &agentConfig,
		logger:     logger,
		caCert:     ca,
		serverName: agentConfig.ServerName,
	}, nil
}

func (a *Agent) Run() {
	wg := sync.WaitGroup{}
	for _, channelConfig := range a.config.Tunnels {
		logger := a.logger.With(slog.String("channel", channelConfig.Name))
		wg.Add(1)
		channel := NewTunnel(channelConfig, a.caCert, a.serverName, logger)
		go channel.Run(&wg)
	}

	wg.Wait()
}
