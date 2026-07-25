package session

import (
	"io"
	"net"
	"testing"

	"go.uber.org/mock/gomock"

	dispatchermocks "github.com/ajsqr/wombat/dispatcher/mocks"
	"github.com/ajsqr/wombat/frame"
)

func TestReceiveSuccess(t *testing.T) {
	s := &Session{
		inbox: make(chan *frame.Frame, 1),
	}

	f := &frame.Frame{
		Ident: frame.DataFrame,
	}

	if err := s.Receive(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := <-s.inbox

	if got != f {
		t.Fatal("frame mismatch")
	}
}

func TestStreamDispatchesFrame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	disp := dispatchermocks.NewMockDispatcher(ctrl)

	client, server := net.Pipe()
	defer client.Close()

	disp.EXPECT().
		Dispatch(gomock.Any()).
		Do(func(f *frame.Frame) {
			if f.ConnectionID != 42 {
				t.Fatalf("wrong connection id")
			}

			if f.Ident != frame.DataFrame {
				t.Fatalf("wrong identifier")
			}

			if string(f.Payload) != "hello" {
				t.Fatalf("unexpected payload %q", string(f.Payload))
			}
		})

	s := &Session{
		sessionID:  42,
		inbox:      make(chan *frame.Frame),
		conn:       server,
		errChan:    make(chan error, 1),
		closed:     make(chan struct{}, 1),
		bufferSize: 1024,
		disp:       disp,
	}

	done := make(chan error)

	go func() {
		done <- s.Stream()
	}()

	_, err := client.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	client.Close()

	err = <-done
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecvLoopWritesPayload(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()

	s := &Session{
		conn:    server,
		inbox:   make(chan *frame.Frame, 1),
		errChan: make(chan error, 1),
	}

	go s.recvLoop()

	err := s.Receive(&frame.Frame{
		Payload: []byte("hello"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 5)

	_, err = io.ReadFull(client, buf)
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != "hello" {
		t.Fatalf("expected hello got %q", string(buf))
	}

	close(s.inbox)
}

func TestLastPayloadBeforeEOFDispatched(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	disp := dispatchermocks.NewMockDispatcher(ctrl)

	client, server := net.Pipe()

	disp.EXPECT().
		Dispatch(gomock.Any()).
		Do(func(f *frame.Frame) {
			if string(f.Payload) != "final" {
				t.Fatalf("expected final payload got %q", string(f.Payload))
			}
		})

	s := &Session{
		sessionID:  1,
		inbox:      make(chan *frame.Frame),
		conn:       server,
		errChan:    make(chan error, 1),
		closed:     make(chan struct{}, 1),
		bufferSize: 1024,
		disp:       disp,
	}

	done := make(chan error)

	go func() {
		done <- s.Stream()
	}()

	_, err := client.Write([]byte("final"))
	if err != nil {
		t.Fatal(err)
	}

	client.Close()

	err = <-done
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMultipleFramesDispatched(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	disp := dispatchermocks.NewMockDispatcher(ctrl)

	client, server := net.Pipe()

	payloads := []string{
		"one",
		"two",
		"three",
	}

	for _, p := range payloads {
		expected := p

		disp.EXPECT().
			Dispatch(gomock.Any()).
			Do(func(f *frame.Frame) {
				if string(f.Payload) != expected {
					t.Fatalf("expected %q got %q", expected, string(f.Payload))
				}
			})
	}

	s := &Session{
		sessionID:  99,
		inbox:      make(chan *frame.Frame),
		conn:       server,
		errChan:    make(chan error, 1),
		closed:     make(chan struct{}, 1),
		bufferSize: 1024,
		disp:       disp,
	}

	done := make(chan error)

	go func() {
		done <- s.Stream()
	}()

	for _, p := range payloads {
		_, err := client.Write([]byte(p))
		if err != nil {
			t.Fatal(err)
		}
	}

	client.Close()

	err := <-done
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
