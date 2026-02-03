package common

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
)

// Frame helpers

// WriteFrame writes [streamID][length][payload]
func WriteFrame(w io.Writer, streamID byte, payload []byte) error {
	length := uint32(len(payload))

	// write streamID
	if _, err := w.Write([]byte{streamID}); err != nil {
		return err
	}

	// write length
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}

	// write payload
	_, err := w.Write(payload)
	return err
}

// ReadFrame reads one frame [streamID][length][payload]
func ReadFrame(r io.Reader) (byte, []byte, error) {
	streamIDBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, streamIDBuf); err != nil {
		return 0, nil, err
	}
	streamID := streamIDBuf[0]

	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return 0, nil, err
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, nil, err
	}

	return streamID, buf, nil
}

// Compress payload
func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

// Decompress payload
func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// WriteFrameCompressed writes [streamID][length][compressed payload]
func WriteFrameCompressed(w io.Writer, streamID byte, payload []byte) error {
	compressed, err := Compress(payload)
	if err != nil {
		return err
	}
	return WriteFrame(w, streamID, compressed)
}

// ReadFrameCompressed reads one frame and decompresses payload
func ReadFrameCompressed(r io.Reader) (byte, []byte, error) {
	streamID, compressed, err := ReadFrame(r)
	if err != nil {
		return 0, nil, err
	}
	payload, err := Decompress(compressed)
	return streamID, payload, err
}
