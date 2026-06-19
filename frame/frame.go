package frame

type Frame struct {
	// ident is an identifier that defines the type of a frame.
	// It can signal opening / closing a new connection, data frame, ping frame etc
	ident FrameIdentifer

	// connectionID defines the client connectionID. This field is used to mux
	// stream of data from multiple clients through a single connection to the agent
	connectionID uint32

	// payloadSize defines the size of the payload present in the frame
	payloadSize uint32

	// payload is an arbitrarly long payload
	payload []byte
}

type FrameIdentifer uint32

const (
	OpenConnection FrameIdentifer = iota
	CloseConnection
	Ping
	Pong
	DataFrame
)
