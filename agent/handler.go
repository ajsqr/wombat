package agent

import (
	"errors"
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
	ah.logger.Info("received a control frame", slog.Any("frameIdentifier", f.Ident), slog.Any("connectionId", f.ConnectionID))
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
		go ah.runSession(s, d)
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

func (ah *AgentHandler) runSession(s *session.Session, disp dispatcher.Dispatcher) {
	err := s.Stream()
	if err == nil {
		// stream ended without any errors
		// likely because it was closed because of a close-connection-frame from server

	} else if errors.Is(err, net.ErrClosed) {
		// stream ended because the connection closed.
		// proceed with removing the session from store
	} else {
		// stream ended because of an unexpected reason
		// teardown required
		ah.logger.Error("session streaming failed", slog.Any("error", err), slog.Any("sessionID", s.GetID()))
		if err := ah.destroySession(s); err != nil {
			ah.logger.Error("unexpected error while destroying session", slog.Any("error", err))
		}
	}

	// stream ended - we must instruct the other end to close the session
	closeFrame := frame.Frame{
		ConnectionID: s.GetID(),
		Ident:        frame.CloseConnection,
	}

	disp.Dispatch(&closeFrame)
	ah.store.Destroy(s.GetID())
	ah.logger.Info("session life-cycle ended", slog.Any("sessionID", s.GetID()))
}

func (ah *AgentHandler) destroySession(s *session.Session) error {
	ah.logger.Info("asked to close session", slog.Any("sessionID", s.GetID()))
	err := s.Close()
	if err != nil {
		return err
	}

	return nil
}
