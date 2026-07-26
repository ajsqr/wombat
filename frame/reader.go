package frame

import (
	"encoding/binary"
	"io"
)

// reader assists with reading a stream of data as frames.
// reader will offer capability to read a stream of data and then unmarshall them as frames.
// reader is not concerned with multiplexing the frames and its only goal is to parse a stream of bytes and unmarhsall them to individual frames

var _ Reader = &FrameReader{}

func NewReader(reader io.Reader) *FrameReader {
	return &FrameReader{
		reader: reader,
	}
}

type Reader interface {
	ReadFrame() (*Frame, error)
}

type FrameReader struct {
	reader     io.Reader
	BufferSize uint32
}

func (r *FrameReader) ReadFrame() (*Frame, error) {
	var err error
	var frame Frame
	frame.Ident, err = r.readFrameIdentifier()
	if err != nil {
		return nil, err
	}

	frame.ConnectionID, err = r.readConnectionId()
	if err != nil {
		return nil, err
	}

	frame.PayloadSize, err = r.readPayloadSize()
	if err != nil {
		return nil, err
	}

	frame.Payload, err = r.readPayload(frame.PayloadSize)
	if err != nil {
		return nil, err
	}

	return &frame, nil

}

func (r *FrameReader) readFrameIdentifier() (FrameIdentifer, error) {
	var fid FrameIdentifer
	fidBytes := make([]byte, 4)
	_, err := io.ReadFull(r.reader, fidBytes)
	if err != nil {
		return fid, err
	}

	fid = FrameIdentifer(binary.LittleEndian.Uint32(fidBytes))
	switch fid {
	case Authentication, OpenConnection, CloseConnection, Ping, Pong, DataFrame:
	default:
		return fid, ErrInvalidFrameIdentifier
	}

	return fid, nil
}

func (r *FrameReader) readConnectionId() (uint32, error) {
	var connID uint32
	buffer := make([]byte, 4)
	_, err := io.ReadFull(r.reader, buffer)
	if err != nil {
		return connID, err
	}

	connID = binary.LittleEndian.Uint32(buffer)
	return connID, nil
}

func (r *FrameReader) readPayloadSize() (uint32, error) {
	var sz uint32
	buffer := make([]byte, 4)
	_, err := io.ReadFull(r.reader, buffer)
	if err != nil {
		return sz, err
	}

	sz = binary.LittleEndian.Uint32(buffer)
	return sz, nil
}

func (r *FrameReader) readPayload(payloadSize uint32) ([]byte, error) {
	payload := make([]byte, payloadSize)
	_, err := io.ReadFull(r.reader, payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}
