package frame

import (
	"bytes"
	"testing"
)

func TestFrameMarshalUnmarshal(t *testing.T) {
	original := Frame{
		ident:        DataFrame,
		connectionID: 42,
		payloadSize:  uint32(len([]byte("hello world"))),
		payload:      []byte("hello world"),
	}

	var buf bytes.Buffer

	writer := FrameWriter{
		writer: &buf,
	}

	err := writer.WriteFrame(original)
	if err != nil {
		t.Fatalf("failed to marshal frame: %v", err)
	}

	reader := FrameReader{
		reader: &buf,
	}

	decoded, err := reader.ReadFrame()
	if err != nil {
		t.Fatalf("failed to unmarshal frame: %v", err)
	}

	if decoded.ident != original.ident {
		t.Fatalf(
			"identifier mismatch. expected=%v got=%v",
			original.ident,
			decoded.ident,
		)
	}

	if decoded.connectionID != original.connectionID {
		t.Fatalf(
			"connectionID mismatch. expected=%d got=%d",
			original.connectionID,
			decoded.connectionID,
		)
	}

	if decoded.payloadSize != original.payloadSize {
		t.Fatalf(
			"payloadSize mismatch. expected=%d got=%d",
			original.payloadSize,
			decoded.payloadSize,
		)
	}

	if !bytes.Equal(decoded.payload, original.payload) {
		t.Fatalf(
			"payload mismatch. expected=%q got=%q",
			original.payload,
			decoded.payload,
		)
	}
}
