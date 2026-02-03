package common

import (
	"encoding/binary"
	"io"
)

type Frame struct {
	StreamID byte
	Payload  []byte
}

// WriteFrame writes a multiplexed frame: [streamID][length][payload]
func WriteFrame(w io.Writer, streamID byte, payload []byte) error {
	length := uint32(len(payload))

	if _, err := w.Write([]byte{streamID}); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}

	_, err := w.Write(payload)
	return err
}

// ReadFrame reads one multiplexed frame
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
