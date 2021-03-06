package packets

import (
	"time"
)

type Message struct {
	Username string
	Message  string
	Time     time.Time
	Channel  string
}

func (m *Message) ToBytes() []byte {
	return dataToBytes(m)
}
func (m *Message) CreatePacket() *Packet {
	return &Packet{
		Header: header{MsgType: 2},
		Data:   m,
	}
}
