package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	CreateUser         = 0x01
	CreateUserResponse = 0x02

	GetUserToken         = 0x03
	GetUserTokenResponse = 0x04
)

type Header struct {
	Opcode   byte
	Length   uint32
	CheckSum uint32
}
type Message struct {
	Header Header
	Body   []byte
}

func calculateChecksum(data []byte) uint32 {
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum
}

func ParseMessage(conn net.Conn) (*Message, error) {
	headerBuf := make([]byte, 9)

	if _, err := io.ReadFull(conn, headerBuf); err != nil {
		return nil, fmt.Errorf("read header failed: %w", err)
	}

	header := Header{
		Opcode:   headerBuf[0],
		Length:   binary.BigEndian.Uint32(headerBuf[1:5]),
		CheckSum: binary.BigEndian.Uint32(headerBuf[5:9]),
	}

	body := make([]byte, header.Length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if calc := calculateChecksum(body); calc != header.CheckSum {
		return nil, fmt.Errorf("checksum mismatch: expected %d, got %d", header.CheckSum, calc)
	}

	return &Message{
		Header: header,
		Body:   body,
	}, nil
}

func EncodeMessage(msg *Message) ([]byte, error) {
	msg.Header.Length = uint32(len(msg.Body))
	msg.Header.CheckSum = calculateChecksum(msg.Body)

	buf := make([]byte, 9+len(msg.Body))
	buf[0] = msg.Header.Opcode
	binary.BigEndian.PutUint32(buf[1:5], msg.Header.Length)
	binary.BigEndian.PutUint32(buf[5:9], msg.Header.CheckSum)
	copy(buf[9:], msg.Body)

	return buf, nil
}
