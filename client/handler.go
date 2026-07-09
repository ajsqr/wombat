package client

import (
	"log/slog"
	"net"

	"github.com/ajsqr/wombat/dispatcher"
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
	// store is where sessions are kept
	store *session.SessionStore
	// hostAddr points to the location in which locally running host is running
	hostAddr string
	logger   *slog.Logger
}

// Handle defines how an agent handles control frames it receives.
// An agent is only supposed to handle Open, Close control frames
// Even though the protocol supports provision for Ping and Pong
// It is currently not implemented. Once the agent receives an Open
// control frame, it needs to generate a new session and
func (ah *AgentHandler) Handle(f *frame.Frame, d dispatcher.Dispatcher) error {
	switch f.Ident {
	case frame.OpenConnection:
		// a new control frame to open a new connection
		// upon receiving this frame, a client should open a new connection with the local server
		// the session should be persisted in the local storage
		conn, err := net.Dial("tcp", ah.hostAddr)
		if err != nil {
			return err
		}

		s := ah.createSession(f.ConnectionID, conn, d)
		go ah.runSession(s)
		return nil
	case frame.CloseConnection:
		s, err := ah.store.GetByID(f.ConnectionID)
		if err != nil {
			ah.logger.Error("attempting to handle close-connection frame", slog.Any("error", err))
			return nil
		}

		err = ah.destroySession(s)
		if err != nil {
			ah.logger.Error("attempting to destroy session", slog.Any("error", err))
			return nil
		}

		return nil
	default:
		return nil
	}

}

func (ah *AgentHandler) createSession(id uint32, conn net.Conn, disp dispatcher.Dispatcher) *session.Session {
	s := session.NewSession(id, conn, disp)
	ah.store.Create(s)
	return s
}

func (ah *AgentHandler) runSession(s *session.Session) {
	err := s.Stream()
	if err != nil {
		// if session streaming is broken we shouldnt error out
		// this must be logged and handler should proceed by closing the session
		ah.logger.Error("failed to stream", slog.Any("error", err))
	}

	// close session
	err = ah.destroySession(s)
	if err != nil {
		ah.logger.Error("failed to destroy session", slog.Any("error", err))
	}

}

func (ah *AgentHandler) destroySession(s *session.Session) error {
	err := s.Close()
	if err != nil {
		ah.logger.Error("failed to close session", slog.Any("error", err))
	}

	ah.store.Destroy(s.GetID())
	return nil
}
