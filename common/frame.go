package common

import (
	"encoding/binary"
	"io"
)

// WriteFrame writes [len][payload]
func WriteFrame(w io.Writer, data []byte) error {
	length := uint32(len(data))

	// write length
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}

	// write payload
	_, err := w.Write(data)
	return err
}

// ReadFrame reads exactly one framed message
func ReadFrame(r io.Reader) ([]byte, error) {
	var length uint32

	// read length
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	buf := make([]byte, length)

	// read payload fully
	_, err := io.ReadFull(r, buf)
	return buf, err
}
