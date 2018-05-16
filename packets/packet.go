package packets

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

type IPacket interface {
	ToBytes() []byte
}
type header struct {
	Size uint16
}
type Packet struct {
	header  header
	MsgType uint16
	Data    IPacket
}

func packetToBytes(p IPacket) []byte {
	data := new(bytes.Buffer)
	enc := gob.NewEncoder(data)
	enc.Encode(p)
	return data.Bytes()
}
func (h *header) ToBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h.Size)
	return buf.Bytes()
}
func (p *Packet) ToBytes() []byte {
	dataBytes := packetToBytes(p.Data)
	p.header.Size = uint16(len(dataBytes))
	fmt.Printf("Packet data size: %d", p.header.Size)
	return append(p.header.ToBytes(), dataBytes...)
}
func (p *Packet) ToBytesFast() []byte {
	dataBytes := packetToBytes(p.Data)
	buff := make([]byte, 2+len(dataBytes))
	binary.BigEndian.PutUint16(buff, p.header.Size)
	copy(buff[2:], dataBytes)
	return buff
}
