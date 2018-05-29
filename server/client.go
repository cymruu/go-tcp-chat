package server

import (
	"chat2/packets"
	"net"
)

type Client struct {
	Username string
	Token    string
	conn     net.Conn
}

func (c *Client) send(packet packets.IPacket) {
	c.conn.Write(packet.ToBytes())
}
func (c *Client) sendData(packet packets.IPacketData) {
	c.conn.Write(packet.CreatePacket().ToBytes())
}
