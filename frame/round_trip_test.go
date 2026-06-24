package frame

import (
	"bytes"
	"testing"
)

func TestFrameMarshalUnmarshal(t *testing.T) {
	original := Frame{
		Ident:        DataFrame,
		ConnectionID: 42,
		PayloadSize:  uint32(len([]byte("hello world"))),
		Payload:      []byte("hello world"),
	}

	var buf bytes.Buffer

	writer := FrameWriter{
		writer: &buf,
	}

	err := writer.WriteFrame(&original)
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

	if decoded.Ident != original.Ident {
		t.Fatalf(
			"Identifier mismatch. expected=%v got=%v",
			original.Ident,
			decoded.Ident,
		)
	}

	if decoded.ConnectionID != original.ConnectionID {
		t.Fatalf(
			"ConnectionID mismatch. expected=%d got=%d",
			original.ConnectionID,
			decoded.ConnectionID,
		)
	}

	if decoded.PayloadSize != original.PayloadSize {
		t.Fatalf(
			"PayloadSize mismatch. expected=%d got=%d",
			original.PayloadSize,
			decoded.PayloadSize,
		)
	}

	if !bytes.Equal(decoded.Payload, original.Payload) {
		t.Fatalf(
			"Payload mismatch. expected=%q got=%q",
			original.Payload,
			decoded.Payload,
		)
	}
}
