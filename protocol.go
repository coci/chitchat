package main

import (
	"encoding/binary"
	"fmt"
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
	if len(data) < 10 {
		return Frame{}, fmt.Errorf("data too short for header: %d bytes", len(data))
	}

	h := Header{
		Magic:       binary.BigEndian.Uint16(data[0:2]),  // 0xD00D
		Version:     data[2],                             // 1
		MessageType: data[3],                             // 15
		StreamId:    binary.BigEndian.Uint32(data[4:8]),  // 10
		Length:      binary.BigEndian.Uint16(data[8:10]), // 20
	}

	body := []byte{}
	if len(data) > 10 {
		body = data[10:]
	}

	return Frame{
		Header:   h,
		LongTail: LongTail{},
		Body:     body,
	}, nil
}
