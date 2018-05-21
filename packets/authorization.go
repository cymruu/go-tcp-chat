package packets

type Authorization struct {
	Username string
	Token    string
}

func (a *Authorization) ToBytes() []byte {
	return dataToBytes(a)
}
func (a *Authorization) CreatePacket() *Packet {
	return &Packet{
		Header: header{MsgType: 1},
		Data:   a,
	}
}
