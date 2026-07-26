package agent

import (
	"fmt"
	"os"

	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/frame"
)

func (t *Tunnel) handshake(fw *frame.FrameWriter, c *config.AgentTunnelConfig) error {
	token := os.Getenv(c.TokenName)
	if token == "" {
		return fmt.Errorf("token not set")
	}

	tokenPayload := []byte(token)

	loginFrame := frame.Frame{
		Ident:       frame.Authentication,
		Payload:     tokenPayload,
		PayloadSize: uint32(len(tokenPayload)),
	}

	return fw.WriteFrame(&loginFrame)
}
