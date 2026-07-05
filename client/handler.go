package client

import (
	"log/slog"
	"net"

	"github.com/ajsqr/wombat/endpoint"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

var _ endpoint.Endpoint = &AgentHandler{}

func NewClientHandler(store *session.SessionStore, hostAddr string, logger *slog.Logger) *AgentHandler {
	return &AgentHandler{
		store:    store,
		hostAddr: hostAddr,
		logger:   logger,
	}
}

type AgentHandler struct {
	store    *session.SessionStore
	hostAddr string
	logger   *slog.Logger
}

// Handle defines how an agent handles control frames it receives.
// An agent is only supposed to handle Open, Close control frames
// Even though the protocol supports provision for Ping and Pong
// It is currently not implemented. Once the agent receives an Open
// control frame, it needs to generate a new session and
func (ah *AgentHandler) Handle(f *frame.Frame) error {
	switch f.Ident {
	case frame.OpenConnection:
		// a new control frame to open a new connection
		// upon receiving this frame, a client should open a new connection with the local server
		conn, err := net.Dial("tcp", ah.hostAddr)
		if err != nil {
			return err
		}

		_ = ah.store.New(conn)
		return nil
	case frame.CloseConnection:
		ah.store.Destroy(f.ConnectionID)
		return nil
	default:
		return nil
	}

}

func (ah *AgentHandler) stream(s *session.Session) {
	err := s.Stream()
	if err != nil {
		ah.logger.Warn("session stream returned an error", slog.Any("error", err))
	}

	err = s.CloseConnection()
	if err != nil {
		ah.logger.Warn("session close", slog.Any("error", err))
	}
}
