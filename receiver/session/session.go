package session

import (
	"io"
	"net"

	"github.com/ajsqr/wombat/dispatcher"
	"github.com/ajsqr/wombat/frame"
	"github.com/ajsqr/wombat/receiver"
)

var _ receiver.Receiver = &Session{}

type Session struct {
	// sessionID is a unique identifier used for routing frames to and from a session
	sessionID uint32
	// a dispatcher delivers frames intented for a session to its inbox
	inbox chan *frame.Frame
	// conn holds the socket connection established with the client
	conn net.Conn
	// errChan holds the error if the session errored out while streaming
	errChan chan error
	// closed signals that a session has been closed
	// this indicates a graceful exit
	closed chan struct{}
	// bufferSize determines the max size of data a session can read from the client
	bufferSize int
	// disp deals with frame transportation between a session and the tunnel
	disp dispatcher.Dispatcher
}

// Receive method exists for receiving frames from a dispatcher
// If the inbox of a session is full, receive returns an error,
// and the caller can terminate the session
func (s *Session) Receive(f *frame.Frame) error {
	select {
	case s.inbox <- f:
		return nil
	default:
		return receiver.ErrCannotReceive
	}
}

func (s *Session) CloseConnection() error {
	return s.conn.Close()
}

func (s *Session) Stream() error {
	go s.sendLoop()
	go s.recvLoop()
	select {
	case err := <-s.errChan:
		return err
	case <-s.closed:
		return nil
	}
}

// sendLoop is responsible for sending data from the client to the dispatcher
func (s *Session) sendLoop() {
	for {
		buffer := make([]byte, s.bufferSize)
		n, err := s.conn.Read(buffer)
		// Callers should always process the n > 0 bytes returned before considering the error err.
		if n > 0 {
			s.dispatchDataFrame(buffer, n)
		}

		switch err {
		case nil:
		case io.EOF:
			// this is a gracefull shutdown
			s.closed <- struct{}{}
			return
		default:
			// an unexpected error occured
			// forward the error to errChan and let handler close the session manually
			s.errChan <- err
			return
		}
	}
}

// recvLoop is responsible for writing the content from frames to the client
func (s *Session) recvLoop() {
	for f := range s.inbox {
		_, err := s.conn.Write(f.Payload)
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

func (s *Session) dispatchDataFrame(buffer []byte, n int) {
	var f frame.Frame
	f.ConnectionID = s.sessionID
	f.Ident = frame.DataFrame
	f.PayloadSize = uint32(n)
	f.Payload = buffer[:n]
	s.disp.Dispatch(&f)
}
