package packets

import (
	"time"
)

type SystemMessage struct {
	Message string
	Time    time.Time
}

func (m *SystemMessage) ToBytes() []byte {
	return dataToBytes(m)
}
func (m *SystemMessage) CreatePacket() *Packet {
	m.Time = time.Now()
	return &Packet{
		Header: header{MsgType: 3},
		Data:   m,
	}
}
