package tunnel

import (
	"net"

	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/endpoint"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver/session"
)

const (
	defaultQueueSize = 100
)

var _ dispatcher.Dispatcher = &Tunnel{}

type Tunnel struct {
	// conn is the tunnel connection
	conn        net.Conn
	frameWriter frame.Writer
	frameReader frame.Reader
	// dispatchQueue queues frames sent from sessions
	dispatchQueue chan *frame.Frame
	// storer has the sessions managed by the tunnel
	storer  *session.SessionStore
	errChan chan error
	handler endpoint.Endpoint
}

func NewTunnel(conn net.Conn, store *session.SessionStore, handler endpoint.Endpoint) *Tunnel {
	tunnel := Tunnel{}

	tunnel.frameWriter = frame.NewWriter(conn)
	tunnel.frameReader = frame.NewReader(conn)
	tunnel.dispatchQueue = make(chan *frame.Frame, defaultQueueSize)
	tunnel.storer = store
	tunnel.errChan = make(chan error, 1)
	tunnel.handler = handler
	return &tunnel
}

func (t *Tunnel) Dispatch(f *frame.Frame) {
	t.dispatchQueue <- f
}

func (t *Tunnel) Stream() error {
	defer t.conn.Close()
	defer close(t.dispatchQueue)
	go t.recvLoop()
	go t.sendLoop()
	return <-t.errChan
}

// sendLoop writs frames in its queue to the tunnel
func (t *Tunnel) sendLoop() {
	for f := range t.dispatchQueue {
		err := t.frameWriter.WriteFrame(f)
		if err != nil {
			t.errChan <- err
			return
		}
	}
}

// recvLoop reads frames from the tunnel and deliveres it to relevant sessions.
// if it receives control frames, it is forwarded to the controller ?
func (t *Tunnel) recvLoop() {
	for {
		f, err := t.frameReader.ReadFrame()
		if err != nil {
			t.errChan <- err
			return
		}

		switch f.Ident {
		case frame.OpenConnection, frame.CloseConnection, frame.Ping, frame.Pong:
			// control frames
			if err := t.handler.Handle(f); err != nil {
				t.errChan <- err
				return
			}
		case frame.DataFrame:
			session, err := t.storer.GetByID(f.ConnectionID)
			if err != nil {
				t.errChan <- err
				return
			}

			err = session.Receive(f)
			if err != nil {
				t.errChan <- err
				return
			}
		default:
			// unknown frame
			// ignored for now
		}
	}
}
