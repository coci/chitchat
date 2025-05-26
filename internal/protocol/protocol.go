package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Header struct {
	Opcode byte
	Length uint32
}

type Message struct {
	Header Header
	Body   []byte
}

func (m *Message) Encode() ([]byte, error) {
	m.Header.Length = uint32(len(m.Body))

	buf := make([]byte, 5+len(m.Body))
	buf[0] = m.Header.Opcode
	binary.BigEndian.PutUint32(buf[1:5], m.Header.Length)
	copy(buf[5:], m.Body)

	return buf, nil
}

func (m *Message) ValidateMessage() error {
	if m.Header.Length == 0 {
		return errors.New("length cannot be zero")
	}

	if m.Header.Length != uint32(len(m.Body)) {
		return errors.New("length of message not equal to length in header")
	}
	return nil
}

func (m *Message) ParseBody() []string {
	parts := bytes.Split(m.Body, []byte(","))

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		result = append(result, string(part))
	}

	return result
}

func ParseMessage(reader io.Reader) (*Message, error) {
	headerBuf := make([]byte, 5)

	if _, err := io.ReadFull(reader, headerBuf); err != nil {
		return nil, fmt.Errorf("read header failed: %w", err)
	}

	header := Header{
		Opcode: headerBuf[0],
		Length: binary.BigEndian.Uint32(headerBuf[1:5]),
	}

	body := make([]byte, header.Length)
	if _, err := io.ReadFull(reader, body); err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return &Message{
		Header: header,
		Body:   body,
	}, nil
}
