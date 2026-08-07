package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ajsqr/wombat/auth"
	envauther "github.com/ajsqr/wombat/auth/env"
	"github.com/ajsqr/wombat/config"
	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

const backOff = 5

type Tunnel struct {
	cert    tls.Certificate
	config  *config.ServerTunnelConfig
	store   *session.SessionStore
	counter atomic.Uint32
	auther  auth.Authenticator
	logger  *slog.Logger
}

func NewTunnel(cfg *config.ServerTunnelConfig, cert tls.Certificate, logger *slog.Logger) *Tunnel {
	var s Tunnel
	s.config = cfg
	s.store = session.NewSessionStore()
	s.counter = atomic.Uint32{}
	s.logger = logger
	s.auther = envauther.NewEnvAuthenticator()
	s.cert = cert
	return &s
}

func (t *Tunnel) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	cfg := &tls.Config{
		Certificates: []tls.Certificate{t.cert},
		MinVersion:   tls.VersionTLS13,
	}

	listener, err := tls.Listen("tcp", t.config.Tunnel, cfg)
	if err != nil {
		t.logger.Error("attempting to create a listener for the tunnel", slog.Any("error", err))
		return
	}

	defer listener.Close()
	for {
		err := t.run(listener)
		if err != nil {
			t.logger.Error("tunnel failed", slog.Any("error", err))
		}

		time.Sleep(time.Second * backOff)
	}

}

func (t *Tunnel) run(listener net.Listener) error {
	handler := NewServerHandler(t.store, t.logger)
	// start tunnel listener
	t.logger.Info("creating tunnel connection")
	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("attempting to accept tunnel connection : %w", err)
	}

	frameWriter := frame.NewWriter(conn)
	frameReader := frame.NewReader(conn)

	t.logger.Info("authenticating tunnel connection")
	err = t.handshake(frameReader, t.config)
	if err != nil {
		conn.Close()
		return fmt.Errorf("error during handshake : %w", err)
	}

	t.logger.Info("successfully autenticated tunnel")
	tunnel := tunnel.NewTunnel(conn, frameWriter, frameReader, t.store, handler)
	t.logger.Info("successfully established a tunnel")

	publicListener, err := net.Listen("tcp", t.config.Public)
	if err != nil {
		return fmt.Errorf("attempting to start listener : %w", err)
	}

	defer publicListener.Close()
	go t.acceptConections(publicListener, tunnel, t.store)
	err = tunnel.Stream()
	if err != nil {
		return fmt.Errorf("tunnel streaming failed : %w", err)
	}

	return nil
}

func (t *Tunnel) acceptConections(listener net.Listener, disp dispatcher.Dispatcher, store *session.SessionStore) {
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
	if err == nil {
		// stream ended without any errors
		// likely because it was closed because of a close-connection-frame
	} else if errors.Is(err, net.ErrClosed) {
		// stream ended because the connection closed.
		// proceed with removing the session from store
	} else {
		// stream ended because of an unexpected reason
		// teardown required
		t.logger.Error("session streaming failed", slog.Any("error", err), slog.Any("sessionID", sess.GetID()))
		if err := t.destroySession(sess); err != nil {
			t.logger.Error("unexpected error while destroying session", slog.Any("error", err))
		}
	}

	if sess.NotifyTunnelOnClose {
		// stream ended - we must instruct the other end to close the session
		closeFrame := frame.Frame{
			ConnectionID: sess.GetID(),
			Ident:        frame.CloseConnection,
		}

		disp.Dispatch(&closeFrame)
		t.logger.Info("dispatched a CloseConnection frame", slog.Any("sessionID", sess.GetID()))
	}

	t.store.Destroy(sess.GetID())

}

func (t *Tunnel) destroySession(sess *session.Session) error {
	t.logger.Info("asked to close session", slog.Any("sessionID", sess.GetID()))
	err := sess.Close()
	if err != nil {
		t.logger.Error("failed to close session", slog.Any("error", err))
	}

	return nil
}
