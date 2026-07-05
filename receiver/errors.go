package receiver

import "errors"

var (
	ErrCannotReceive = errors.New("unable to receive as the channel is full")

	ErrSessionNotFound = errors.New("session not found")
)
