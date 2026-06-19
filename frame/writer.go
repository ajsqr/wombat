package frame

import (
	"encoding/binary"
	"io"
)

// writer deals with marshalling a frame to a slice of bytes so that it can be sent through the stream

type Writer interface {
	WriteFrame(f Frame) error
}

type FrameWriter struct {
	writer io.Writer
}

func (fw *FrameWriter) WriteFrame(f Frame) error {
	err := fw.writeFrameIdentifier(f.ident)
	if err != nil {
		return err
	}

	err = fw.writeConnectionID(f.connectionID)
	if err != nil {
		return err
	}

	err = fw.writePayloadSize(f.payloadSize)
	if err != nil {
		return err
	}

	err = fw.writePayload(f.payload)
	if err != nil {
		return err
	}

	return nil
}

func (fw *FrameWriter) writeFrameIdentifier(ident FrameIdentifer) error {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(ident))
	_, err := fw.writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (fw *FrameWriter) writeConnectionID(connID uint32) error {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, connID)
	_, err := fw.writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (fw *FrameWriter) writePayloadSize(sz uint32) error {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, sz)
	_, err := fw.writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (fw *FrameWriter) writePayload(data []byte) error {
	_, err := fw.writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}
