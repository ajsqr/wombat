package server

import (
	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/frame"
)

func (t *Tunnel) handshake(fr *frame.FrameReader, c *config.ServerTunnelConfig) error {
	f, err := fr.ReadFrame()
	if err != nil {
		return err
	}

	return t.auther.Authenticate(c.TokenName, f.Payload)
}
