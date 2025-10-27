package protocol

import (
	"chitchat/pkg/logger"
	"encoding/binary"
	"fmt"
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
	LongTail *LongTail
	Body     []byte
}
type Gossip struct {
	IProtocol

	Logger logger.ILogger
}

func (g Gossip) ParseMessage(data []byte) (IFrame, error) {
	if len(data) < 10 {
		return Frame{}, fmt.Errorf("data too short for header: %d bytes", len(data))
	}

	h := Header{
		Magic:       binary.BigEndian.Uint16(data[0:2]),
		Version:     data[2],
		MessageType: data[3],
		StreamId:    binary.BigEndian.Uint32(data[4:8]),
		Length:      binary.BigEndian.Uint16(data[8:10]),
	}

	body := []byte{}
	if len(data) > 10 {
		body = data[10:]
	}

	return Frame{
		Header: h,
		Body:   body,
	}, nil
}

func (g Gossip) SerializeMessage(f IFrame) []byte {

	frame, ok := f.(Frame)

	if !ok {
		return []byte{}
	}

	isLongHeader := frame.LongTail != nil

	var headerSize int

	if isLongHeader {
		headerSize = 38
	} else {
		headerSize = 10
	}

	data := make([]byte, headerSize+len(frame.Body))

	binary.BigEndian.PutUint16(data[0:2], frame.Header.Magic)
	data[2] = frame.Header.Version
	data[3] = frame.Header.MessageType
	binary.BigEndian.PutUint32(data[4:8], frame.Header.StreamId)
	binary.BigEndian.PutUint16(data[8:10], frame.Header.Length)

	if isLongHeader {
		// Serialize Long Tail (post-auth)
		binary.BigEndian.PutUint32(data[10:14], frame.LongTail.SessionId)
		binary.BigEndian.PutUint64(data[14:22], frame.LongTail.Sequence)
		binary.BigEndian.PutUint64(data[22:38], frame.LongTail.Tag)
		// Copy body
		copy(data[38:], frame.Body)
	} else {
		// Copy body for short header
		copy(data[10:], frame.Body)
	}

	return data
}

func (g Gossip) String(frame IFrame) string {
	f, ok := frame.(Frame)
	if !ok {
		g.Logger.Error("The frame is not a IFrame")
	}

	if f.LongTail != nil {
		return fmt.Sprintf(
			"Frame{Magic: 0x%04X, Version: %d, MsgType: %d, StreamID: %d, Length: %d, SessionID: %d, Sequence: %d, Tag: %d, Body: %q}",
			f.Header.Magic,
			f.Header.Version,
			f.Header.MessageType,
			f.Header.StreamId,
			f.Header.Length,
			f.LongTail.SessionId,
			f.LongTail.Sequence,
			f.LongTail.Tag,
			f.Body,
		)
	}
	return fmt.Sprintf(
		"Frame{Magic: 0x%04X, Version: %d, MsgType: %d, StreamID: %d, Length: %d, Body: %q}",
		f.Header.Magic,
		f.Header.Version,
		f.Header.MessageType,
		f.Header.StreamId,
		f.Header.Length,
		f.Body,
	)
}
