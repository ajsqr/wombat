package receiver

import "github.com/ajsqr/wombat/frame"

type Receiver interface {
	// Receive exposes the capability using which an implementation can receive frames from a sender
	// It will return an error if the receiver is unable to receive the frame.
	Receive(f *frame.Frame) error
}
