package server

import (
	"chat2/packets"
	"fmt"
)

type Channel struct {
	name         string
	participants []*Client
}

func (c *Channel) join(client *Client) {
	c.participants = append(c.participants, client)

	c.broadcast(&packets.SystemMessage{Message: fmt.Sprintf("%s joined channel %s", client.Username, c.name)})
}

func (c *Channel) leave(client *Client) {
	for i, participant := range c.participants {
		if client == participant {
			c.participants = append(c.participants[:i], c.participants[i+1:]...)
			c.broadcast(&packets.SystemMessage{Message: fmt.Sprintf("%s has left channel %s", client.Username, c.name)})
		}
	}
}
func (c *Channel) broadcast(packet packets.IPacketData) {
	for _, c := range c.participants {
		c.sendData(packet)
	}
}
