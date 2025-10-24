package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	Magic1  = 0xD0
	Magic2  = 0x0D
	Version = 1
)

type Header struct {
	Magic       uint16
	Version     uint8
	MessageType uint8
	StreamId    uint32
	Length      uint16
}

type LongTail struct {
	SessionId uint32
	Sequence  uint64
	Tag       uint64
}

type Frame struct {
	Header   Header
	LongTail LongTail
	Body     []byte
}

func EncodeMessage(data []byte) (Frame, error) {
	h := &Header{}
	buf := bytes.NewReader(data)

	binary.Read(buf, binary.BigEndian, &h.Magic)
	binary.Read(buf, binary.BigEndian, &h.Version)
	binary.Read(buf, binary.BigEndian, &h.MessageType)
	binary.Read(buf, binary.BigEndian, &h.StreamId)
	binary.Read(buf, binary.BigEndian, &h.Length)

	return Frame{
		Header:   *h,
		LongTail: LongTail{},
		Body:     data,
	}, nil
}
