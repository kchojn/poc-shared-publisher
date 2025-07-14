package network

import (
	"encoding/binary"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

// Codec handles message encoding/decoding.
type Codec struct {
	maxMessageSize int
}

// NewCodec creates a new codec.
func NewCodec(maxSize int) *Codec {
	return &Codec{
		maxMessageSize: maxSize,
	}
}

// Encode encodes a message with a length prefix.
func (c *Codec) Encode(msg *proto.Message) ([]byte, error) {
	data, err := proto.Marshal(*msg)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	if len(data) > c.maxMessageSize {
		return nil, fmt.Errorf("message too large: %d > %d", len(data), c.maxMessageSize)
	}

	// Length prefix (4 bytes, big endian)
	result := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(result[:4], uint32(len(data))) //nolint:gosec // G115
	copy(result[4:], data)

	return result, nil
}

// Decode decodes a message from a reader.
func (c *Codec) Decode(r io.Reader) (*proto.Message, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if int(length) > c.maxMessageSize {
		return nil, fmt.Errorf("message too large: %d > %d", length, c.maxMessageSize)
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	var msg proto.Message
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &msg, nil
}
