package packets

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
)

const (
	headerSize = 2 + 2
)

type IPacket interface {
	ToBytes() []byte
}
type header struct {
	Size    uint16
	MsgType uint16
}
type Packet struct {
	Header header
	Data   IPacket
}

func (h *header) ToBytes() []byte {
	buff := make([]byte, headerSize)
	binary.BigEndian.PutUint16(buff[0:2], h.Size)
	binary.BigEndian.PutUint16(buff[2:4], h.MsgType)
	return buff
}
func packetToBytes(p IPacket) []byte {
	data := new(bytes.Buffer)
	enc := gob.NewEncoder(data)
	enc.Encode(p)
	return data.Bytes()
}
func ReadHeader(headerBytes []byte) header {
	return header{
		Size:    binary.BigEndian.Uint16(headerBytes[:2]),
		MsgType: binary.BigEndian.Uint16(headerBytes[2:4]),
	}
}
func (p *Packet) ToBytes() []byte {
	buf := new(bytes.Buffer)
	dataBytes := packetToBytes(p.Data)
	p.Header.Size = uint16(len(dataBytes))
	buf.Write(p.Header.ToBytes())
	buf.Write(dataBytes)
	return buf.Bytes()
}
func (p *Packet) ToBytesFast() []byte {
	dataBytes := packetToBytes(p.Data)
	p.Header.Size = uint16(len(dataBytes))
	buff := make([]byte, headerSize+len(dataBytes))
	copy(buff[:headerSize], p.Header.ToBytes())
	copy(buff[headerSize:], dataBytes)
	return buff
}
func FromBytes(header header, databytes []byte) (Packet, error) {
	dec := gob.NewDecoder(bytes.NewReader(databytes))
	p := Packet{Header: header}
	switch header.MsgType {
	case 1:
		var received Authorization
		err := dec.Decode(&received)
		if err != nil {
			fmt.Print(err.Error())
		}
		p.Data = &received
	case 2:
		var received Message
		err := dec.Decode(&received)
		if err != nil {
			fmt.Print(err.Error())
		}
		p.Data = &received
	default:
		return p, errors.New("Corrupted packet")
	}
	return p, nil
}
