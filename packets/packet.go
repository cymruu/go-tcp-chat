package packets

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const (
	headerSize = 2 + 2
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
func ReadHeader(header []byte) (uint16, uint16) {
	return binary.BigEndian.Uint16(header[:2]), binary.BigEndian.Uint16(header[2:4])
}
func (p *Packet) CreateHeader() []byte {
	if p.header.Size == 0 {
		p.header.Size = uint16(len(p.Data.ToBytes()))
	}
	buff := make([]byte, headerSize)
	binary.BigEndian.PutUint16(buff[:2], p.header.Size)
	binary.BigEndian.PutUint16(buff[2:], p.MsgType)
	return buff
}
func (p *Packet) ToBytes() []byte {
	buf := new(bytes.Buffer)
	dataBytes := packetToBytes(p.Data)
	p.header.Size = uint16(len(dataBytes))
	buf.Write(p.CreateHeader())
	buf.Write(dataBytes)
	return buf.Bytes()
}
func (p *Packet) ToBytesFast() []byte {
	dataBytes := packetToBytes(p.Data)
	buff := make([]byte, headerSize+len(dataBytes))
	copy(buff[:4], p.CreateHeader())
	copy(buff[4:], dataBytes)
	return buff
}
