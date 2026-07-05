package frame

import (
	"encoding/binary"
	"io"
)

// writer deals with marshalling a frame to a slice of bytes so that it can be sent through the stream

func NewWriter(writer io.Writer) *FrameWriter {
	return &FrameWriter{
		writer: writer,
	}
}

type Writer interface {
	WriteFrame(f *Frame) error
}

type FrameWriter struct {
	writer io.Writer
}

func (fw *FrameWriter) WriteFrame(f *Frame) error {
	err := fw.writeFrameIdentifier(f.Ident)
	if err != nil {
		return err
	}

	err = fw.writeConnectionID(f.ConnectionID)
	if err != nil {
		return err
	}

	err = fw.writePayloadSize(f.PayloadSize)
	if err != nil {
		return err
	}

	err = fw.writePayload(f.Payload)
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
