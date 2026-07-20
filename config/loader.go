package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	wombat    = "wombat"
	agentCfg  = "agent-config.json"
	serverCfg = "server-config.json"
	logs      = "logs.json"
)

func LoadAgentConfig(config any) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configFilePath := agentConfigPath(configDir)

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, config)
	if err != nil {
		return err
	}

	return nil
}

func LoadServerConfig(config any) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configFilePath := serverConfigPath(configDir)

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, config)
	if err != nil {
		return err
	}

	return nil
}

func agentConfigPath(rootConfig string) string {
	return filepath.Join(rootConfig, wombat, agentCfg)
}

func serverConfigPath(rootConfig string) string {
	return filepath.Join(rootConfig, wombat, serverCfg)
}
