package protocol

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

type IProtocol interface {
	ParseMessage(data []byte) (Frame, error)
	SerializeMessage(Frame) []byte
}
