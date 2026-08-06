package server

import (
	"errors"
	"log/slog"
	"net"

	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/endpoint"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver"
	"github.com/ajsqr/wombat/receiver/session"
)

var _ endpoint.Endpoint = &ServersHandler{}

type ServersHandler struct {
	// store is where sessions are kept
	store *session.SessionStore
	// hostAddr points to the location in which locally running host is running
	logger *slog.Logger
}

func NewServerHandler(store *session.SessionStore, logger *slog.Logger) *ServersHandler {
	return &ServersHandler{
		store:  store,
		logger: logger,
	}
}

func (sh *ServersHandler) Handle(f *frame.Frame, disp dispatcher.Dispatcher) error {
	switch f.Ident {
	case frame.OpenConnection:
		// a server should not be receiving open frames
		// ignore for now
		return nil
	case frame.CloseConnection:
		// client closed the connection
		s, err := sh.store.GetByID(f.ConnectionID)
		switch err {
		case nil:
		case receiver.ErrSessionNotFound:
			sh.logger.Warn("attempting to handle close-connection frame", slog.Any("error", err))
			return nil
		default:
			sh.logger.Error("attempting to handle close-connection frame", slog.Any("error", err))
			return nil
		}

		err = sh.destroySession(s)
		if err != nil {
			sh.logger.Error("attempting to destroy session", slog.Any("error", err))
			return nil
		}
	default:
		return nil
	}
	return nil
}

func (sh *ServersHandler) destroySession(s *session.Session) error {
	err := s.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		sh.logger.Error("failed to close session", slog.Any("error", err))
	}

	sh.store.Destroy(s.GetID())
	return nil
}
