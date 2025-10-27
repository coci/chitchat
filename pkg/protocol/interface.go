package protocol

type IFrame interface{}

type IProtocol interface {
	ParseMessage(data []byte) (IFrame, error)
	SerializeMessage(IFrame) []byte
	String(IFrame) string
}
