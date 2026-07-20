package server

import (
	"log/slog"
	"net"
	"sync"
	"sync/atomic"

	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

type Tunnel struct {
	config  *config.ServerTunnelConfig
	store   *session.SessionStore
	counter atomic.Uint32
	logger  *slog.Logger
}

func NewTunnel(cfg *config.ServerTunnelConfig, logger *slog.Logger) *Tunnel {
	var s Tunnel
	s.config = cfg
	s.store = session.NewSessionStore()
	s.counter = atomic.Uint32{}
	s.logger = logger
	return &s
}

func (t *Tunnel) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	handler := NewServerHandler(t.store, t.logger)
	// start tunnel listener
	t.logger.Info("creating tunnel connection")
	listener, err := net.Listen("tcp", t.config.Tunnel)
	if err != nil {
		t.logger.Error("attempting to create a listener for the tunnel", slog.Any("error", err))
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		t.logger.Error("attempting to accept tunnel connection", slog.Any("error", err))
		return
	}

	tunnel := tunnel.NewTunnel(conn, t.store, handler)
	t.logger.Info("successfully established a tunnel")
	go t.acceptConections(tunnel, t.store)
	err = tunnel.Stream()
	if err != nil {
		t.logger.Error("tunnel streaming failed", slog.Any("error", err))
		return
	}
}

func (t *Tunnel) acceptConections(disp dispatcher.Dispatcher, store *session.SessionStore) {
	listener, err := net.Listen("tcp", t.config.Public)
	if err != nil {
		t.logger.Error("attempting to start listener", slog.Any("error", err))
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			t.logger.Error("attempting to accept connection", slog.Any("error", err))
			return
		}

		sessionID := t.counter.Add(1)
		sess := session.NewSession(sessionID, conn, disp)
		store.Create(sess)
		openFrame := frame.Frame{
			Ident:        frame.OpenConnection,
			ConnectionID: sessionID,
		}

		disp.Dispatch(&openFrame)
		go t.runSession(sess, disp)
	}
}

func (t *Tunnel) runSession(sess *session.Session, disp dispatcher.Dispatcher) {
	err := sess.Stream()
	if err != nil {
		// if session streaming is broken we shouldnt error out
		// this must be logged and handler should proceed by closing the session
		t.logger.Error("failed to stream", slog.Any("error", err))
	}

	// stream ended - we must instruct the other end to close the session
	closeFrame := frame.Frame{
		ConnectionID: sess.GetID(),
		Ident:        frame.CloseConnection,
	}

	disp.Dispatch(&closeFrame)
	// close session
	err = t.destroySession(sess)
	if err != nil {
		t.logger.Error("failed to destroy session", slog.Any("error", err))
	}

}

func (t *Tunnel) destroySession(sess *session.Session) error {
	err := sess.Close()
	if err != nil {
		t.logger.Error("failed to close session", slog.Any("error", err))
	}

	t.store.Destroy(sess.GetID())
	return nil
}
