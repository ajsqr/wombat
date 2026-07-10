package server

import (
	"log/slog"
	"net"
	"os"
	"sync/atomic"

	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

type Config struct {
	// TunnelAddr is the ip address of the wombat server. A client initiates the tunnel connection to this address.
	// It should be of the format host:port
	// example 192.168.1.15:1234
	TunnelAddr string `json:"tunnelAddr"`
	ServerAddr string `json:"serverAddr"`
}

type Server struct {
	config  *Config
	store   *session.SessionStore
	counter atomic.Uint32
	logger  *slog.Logger
}

func NewServer(cfg *Config) *Server {
	var s Server
	s.config = cfg
	s.store = session.NewSessionStore()
	s.counter = atomic.Uint32{}
	s.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &s
}

func (s *Server) Run() error {
	handler := NewServerHandler(s.store, s.logger)
	// start tunnel listener
	s.logger.Info("creating tunnel connection")
	listener, err := net.Listen("tcp", s.config.TunnelAddr)
	if err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	tunnel := tunnel.NewTunnel(conn, s.store, handler)
	s.logger.Info("successfully established a tunnel")
	go s.acceptConections(tunnel, s.store)
	return tunnel.Stream()
}

func (s *Server) acceptConections(disp dispatcher.Dispatcher, store *session.SessionStore) {
	listener, err := net.Listen("tcp", s.config.ServerAddr)
	if err != nil {
		s.logger.Error("attempting to start listener", slog.Any("error", err))
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("attempting to accept connection", slog.Any("error", err))
			return
		}

		sessionID := s.counter.Add(1)
		sess := session.NewSession(sessionID, conn, disp)
		store.Create(sess)
		openFrame := frame.Frame{
			Ident:        frame.OpenConnection,
			ConnectionID: sessionID,
		}

		disp.Dispatch(&openFrame)
		go s.runSession(sess, disp)
	}
}

func (s *Server) runSession(sess *session.Session, disp dispatcher.Dispatcher) {
	err := sess.Stream()
	if err != nil {
		// if session streaming is broken we shouldnt error out
		// this must be logged and handler should proceed by closing the session
		s.logger.Error("failed to stream", slog.Any("error", err))
	}

	// stream ended - we must instruct the other end to close the session
	closeFrame := frame.Frame{
		ConnectionID: sess.GetID(),
		Ident:        frame.CloseConnection,
	}

	disp.Dispatch(&closeFrame)
	// close session
	err = s.destroySession(sess)
	if err != nil {
		s.logger.Error("failed to destroy session", slog.Any("error", err))
	}

}

func (s *Server) destroySession(sess *session.Session) error {
	err := sess.Close()
	if err != nil {
		s.logger.Error("failed to close session", slog.Any("error", err))
	}

	s.store.Destroy(sess.GetID())
	return nil
}
