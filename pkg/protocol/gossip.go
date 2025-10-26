package protocol

import (
	"encoding/binary"
	"fmt"
)

type Gossip struct {
	IProtocol
}

func (g Gossip) ParseMessage(data []byte) (Frame, error) {
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

func (g Gossip) SerializeMessage(f Frame) []byte {

	isLongHeader := f.LongTail != nil

	var headerSize int

	if isLongHeader {
		headerSize = 38
	} else {
		headerSize = 10
	}

	data := make([]byte, headerSize+len(f.Body))

	binary.BigEndian.PutUint16(data[0:2], f.Header.Magic)
	data[2] = f.Header.Version
	data[3] = f.Header.MessageType
	binary.BigEndian.PutUint32(data[4:8], f.Header.StreamId)
	binary.BigEndian.PutUint16(data[8:10], f.Header.Length)

	if isLongHeader {
		// Serialize Long Tail (post-auth)
		binary.BigEndian.PutUint32(data[10:14], f.LongTail.SessionId)
		binary.BigEndian.PutUint64(data[14:22], f.LongTail.Sequence)
		binary.BigEndian.PutUint64(data[22:38], f.LongTail.Tag)
		// Copy body
		copy(data[38:], f.Body)
	} else {
		// Copy body for short header
		copy(data[10:], f.Body)
	}

	return data
}

func (g Gossip) FrameToString(f Frame) string {
	return fmt.Sprintf(
		"GOSSIP HEADER {\n"+
			"  Magic: 0x%04X\n"+
			"  Version: %d\n"+
			"  MsgType: %d\n"+
			"  StreamID: %d\n"+
			"  Length: %d\n"+
			"  Payload: %q\n"+
			"}",
		f.Header.Magic,
		f.Header.Version,
		f.Header.MessageType,
		f.Header.StreamId,
		f.Header.Length,
		f.Body,
	)
}
