package packets

type Authorization struct {
	Username string
	Token    string
}

func (a *Authorization) ToBytes() []byte {
	return packetToBytes(a)
}
func (a *Authorization) CreatePacket() *Packet {
	return &Packet{
		header: header{MsgType: 1},
		Data:   a,
	}
}
