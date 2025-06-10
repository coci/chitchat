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

func (m *Message) Validate() error {
	if m.Header.Length == 0 {
		return errors.New("message length cannot be zero")
	}
	if m.Header.Length != uint32(len(m.Body)) {
		return errors.New("header length does not match body length")
	}
	return nil
}

func (m *Message) ParseBody() []string {
	parts := bytes.Split(m.Body, []byte(","))
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = string(part)
	}
	return result
}

func ParseMessage(r io.Reader) (*Message, error) {
	headerBuf := make([]byte, 5)
	if _, err := io.ReadFull(r, headerBuf); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	header := Header{
		Opcode: headerBuf[0],
		Length: binary.BigEndian.Uint32(headerBuf[1:5]),
	}
	body := make([]byte, header.Length)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return &Message{Header: header, Body: body}, nil
}

func NewMessage(opcode byte, body []byte) *Message {
	return &Message{
		Header: Header{
			Opcode: opcode,
			Length: uint32(len(body)),
		},
		Body: body,
	}
}
