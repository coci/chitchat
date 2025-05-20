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

	GetChatroomLists         = 0x05
	GetChatroomListsResponse = 0x06

	CreateChatroom         = 0x07
	CreateChatroomResponse = 0x08

	JoinChatroom         = 0x09
	JoinChatroomResponse = 0x0A

	GetChatroomUser         = 0x0B
	GetChatroomUserResponse = 0x0C

	GetChatroomMessages         = 0x0D
	GetChatroomMessagesResponse = 0x0E

	SendMessageTOChatroom         = 0x0F
	SendMessageTOChatroomResponse = 0x10

	BroadcastMessageToChatroom = 0x11

	BroadCastUserJoinedToChatroom = 0x12

	UserLoggedOut = 0x13

	ErrorResponse = 0xFF
)

type Header struct {
	Opcode byte
	Length uint32
}
type Message struct {
	Header Header
	Body   []byte
}

func ParseMessage(conn net.Conn) (*Message, error) {
	headerBuf := make([]byte, 5)

	if _, err := io.ReadFull(conn, headerBuf); err != nil {
		return nil, fmt.Errorf("read header failed: %w", err)
	}

	header := Header{
		Opcode: headerBuf[0],
		Length: binary.BigEndian.Uint32(headerBuf[1:5]),
	}

	body := make([]byte, header.Length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return &Message{
		Header: header,
		Body:   body,
	}, nil
}

func EncodeMessage(msg *Message) ([]byte, error) {
	msg.Header.Length = uint32(len(msg.Body))

	buf := make([]byte, 5+len(msg.Body))
	buf[0] = msg.Header.Opcode
	binary.BigEndian.PutUint32(buf[1:5], msg.Header.Length)
	copy(buf[5:], msg.Body)

	return buf, nil
}
