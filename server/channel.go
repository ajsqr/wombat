package server

import (
	"log/slog"
	"net"
	"sync"
	"sync/atomic"

	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/dispatcher/tunnel"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

type ChannelConfig struct {
	// TunnelAddr is the ip address of the wombat server. A client initiates the tunnel connection to this address.
	// It should be of the format host:port
	// example 192.168.1.15:1234
	TunnelAddr string `json:"tunnelAddr"`
	ServerAddr string `json:"serverAddr"`
	Name       string `json:"name"`
}

type Channel struct {
	config  *ChannelConfig
	store   *session.SessionStore
	counter atomic.Uint32
	logger  *slog.Logger
}

func NewChannel(cfg *ChannelConfig, logger *slog.Logger) *Channel {
	var s Channel
	s.config = cfg
	s.store = session.NewSessionStore()
	s.counter = atomic.Uint32{}
	s.logger = logger
	return &s
}

func (c *Channel) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	handler := NewServerHandler(c.store, c.logger)
	// start tunnel listener
	c.logger.Info("creating tunnel connection")
	listener, err := net.Listen("tcp", c.config.TunnelAddr)
	if err != nil {
		c.logger.Error("attempting to create a listener for the tunnel", slog.Any("error", err))
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		c.logger.Error("attempting to accept tunnel connection", slog.Any("error", err))
		return
	}

	tunnel := tunnel.NewTunnel(conn, c.store, handler)
	c.logger.Info("successfully established a tunnel")
	go c.acceptConections(tunnel, c.store)
	err = tunnel.Stream()
	if err != nil {
		c.logger.Error("tunnel streaming failed", slog.Any("error", err))
		return
	}
}

func (c *Channel) acceptConections(disp dispatcher.Dispatcher, store *session.SessionStore) {
	listener, err := net.Listen("tcp", c.config.ServerAddr)
	if err != nil {
		c.logger.Error("attempting to start listener", slog.Any("error", err))
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			c.logger.Error("attempting to accept connection", slog.Any("error", err))
			return
		}

		sessionID := c.counter.Add(1)
		sess := session.NewSession(sessionID, conn, disp)
		store.Create(sess)
		openFrame := frame.Frame{
			Ident:        frame.OpenConnection,
			ConnectionID: sessionID,
		}

		disp.Dispatch(&openFrame)
		go c.runSession(sess, disp)
	}
}

func (c *Channel) runSession(sess *session.Session, disp dispatcher.Dispatcher) {
	err := sess.Stream()
	if err != nil {
		// if session streaming is broken we shouldnt error out
		// this must be logged and handler should proceed by closing the session
		c.logger.Error("failed to stream", slog.Any("error", err))
	}

	// stream ended - we must instruct the other end to close the session
	closeFrame := frame.Frame{
		ConnectionID: sess.GetID(),
		Ident:        frame.CloseConnection,
	}

	disp.Dispatch(&closeFrame)
	// close session
	err = c.destroySession(sess)
	if err != nil {
		c.logger.Error("failed to destroy session", slog.Any("error", err))
	}

}

func (c *Channel) destroySession(sess *session.Session) error {
	err := sess.Close()
	if err != nil {
		c.logger.Error("failed to close session", slog.Any("error", err))
	}

	c.store.Destroy(sess.GetID())
	return nil
}
