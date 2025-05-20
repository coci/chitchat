package parser

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
	header Header
	body   []byte
}

func Encoder(msg Message) []byte {
	return []byte{}
}
