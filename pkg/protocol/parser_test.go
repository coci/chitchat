package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParseMessageWithValidMessage(t *testing.T) {
	var buf bytes.Buffer
	body := []byte("test,test123456")

	// write header
	buf.WriteByte(0x01)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(body)))
	buf.Write(lengthBytes)

	// write body
	buf.Write(body)

	msg, err := ParseMessage(&buf)
	if err != nil {
		t.Fatalf("ParseMessage failed: %v", err)
	}

	if msg.Header.Opcode != CreateUser {
		t.Errorf("Opcode: got %x, want %x", msg.Header.Opcode, 0x01)
	}

	if msg.Header.Length != uint32(len(body)) {
		t.Errorf("Length: got %d, want %d", msg.Header.Length, 4)
	}

}

func TestParseMessageWithZeroLength(t *testing.T) {
	var buf bytes.Buffer

	// write header
	buf.WriteByte(0x01)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, 0)
	buf.Write(lengthBytes)

	msg, err := ParseMessage(&buf)
	if err != nil {
		t.Fatalf("ParseMessage failed: %v", err)
	}

	if msg.Header.Opcode != CreateUser {
		t.Errorf("Opcode: got %x, want %x", msg.Header.Opcode, 0x01)
	}

	if msg.Header.Length != 0 {
		t.Errorf("Length: got %d, want %d", msg.Header.Length, 4)
	}

	if !bytes.Equal(msg.Body, []byte{}) {
		t.Errorf("expected body to be empty, but got %v", msg.Body)
	}
}

func TestParseMessageWithIncompleteHeader(t *testing.T) {
	var buf bytes.Buffer

	// write header
	buf.WriteByte(0x01)

	_, err := ParseMessage(&buf)

	if err == nil {
		t.Fatal("expected error for incomplete header, got nil")
	}

}

func TestParseMessageWithIncompleteBody(t *testing.T) {
	var buf bytes.Buffer
	body := []byte("test,test123456")

	// write header
	buf.WriteByte(0x01)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(body))+10)
	buf.Write(lengthBytes)

	// write body
	buf.Write(body)

	_, err := ParseMessage(&buf)

	if err == nil {
		t.Fatal("expected error for incomplete body, got nil")
	}
}

func TestParseMessageWithEmptyReader(t *testing.T) {
	var buf bytes.Buffer

	_, err := ParseMessage(&buf)

	if err == nil {
		t.Fatal("expected error for empty data, got nil")
	}
}

func TestParseMessageWithMultipleMessages(t *testing.T) {
	var buf bytes.Buffer

	writeMessage := func(opcode byte, body []byte) {
		buf.WriteByte(opcode)
		lengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBytes, uint32(len(body)))
		buf.Write(lengthBytes)
		buf.Write(body)
	}

	body1 := []byte("test,test123456")
	body2 := []byte("test1")

	writeMessage(0x01, body1)
	writeMessage(0x02, body2)

	// Parse first message
	msg1, err := ParseMessage(&buf)
	if err != nil {
		t.Fatalf("ParseMessage first message failed: %v", err)
	}
	if msg1.Header.Opcode != 0x01 || !bytes.Equal(msg1.Body, body1) {
		t.Errorf("First message incorrect: got opcode %x body %q", msg1.Header.Opcode, msg1.Body)
	}

	// Parse second message
	msg2, err := ParseMessage(&buf)
	if err != nil {
		t.Fatalf("ParseMessage second message failed: %v", err)
	}
	if msg2.Header.Opcode != 0x02 || !bytes.Equal(msg2.Body, body2) {
		t.Errorf("Second message incorrect: got opcode %x body %q", msg2.Header.Opcode, msg2.Body)
	}

	// After reading two messages, buffer should be empty
	rest := buf.Bytes()
	if len(rest) != 0 {
		t.Errorf("Expected buffer to be empty after reading messages, but has %d bytes left", len(rest))
	}
}

func TestEncodeMessageWithValidMessage(t *testing.T) {
	msg := &Message{
		Header: Header{
			Opcode: CreateUser,
			Length: uint32(len("test,test123456")),
		},
		Body: []byte("test,test123456"),
	}

	byt, err := EncodeMessage(msg)

	if err != nil {
		t.Fatalf("EncodeMessage failed: %v", err)
	}

	if len(byt) != 5+len(msg.Body) {
		t.Errorf("encoded length: got %d, want %d", len(byt), len(byt))
	}

	if byt[0] != msg.Header.Opcode {
		t.Errorf("opcode byte: got %x, want %x", byt[0], msg.Header.Opcode)
	}

	lengthFormEncodeMessage := binary.BigEndian.Uint32(byt[1:5])
	if lengthFormEncodeMessage != msg.Header.Length {
		t.Errorf("length bytes: got %d, want %d", lengthFormEncodeMessage, len(msg.Body))
	}

	if !bytes.Equal(msg.Body, []byte("test,test123456")) {
		t.Errorf("encoded body: got %v, want %v", msg.Body, "test,test123456")
	}
}
