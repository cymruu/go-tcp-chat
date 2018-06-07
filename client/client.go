package main

import (
	"bufio"
	"chat2/packets"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	conn     net.Conn
	reader   *bufio.Reader
	Username string
	recv     chan packets.Packet
}

var activeChannel = "main"

func formatTime(t time.Time) string {
	return t.Format("15:04")
}
func (c *Client) readString() string {
	text, _ := c.reader.ReadString('\n')
	return text[:len(text)-2] //remove last two bytes which are \r\n
}
func startClient() (*Client, error) {
	connection, err := net.Dial("tcp", ":3300")
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &Client{
		conn:   connection,
		reader: bufio.NewReader(os.Stdin),
		recv:   make(chan packets.Packet),
	}, nil
}
func (c *Client) login() {
	username := ""
	for username == "" {
		fmt.Printf("Please enter your username: ")
		username = c.readString()
	}
	c.Username = username
	t, _ := time.Now().MarshalBinary()
	loginPacket := packets.Authorization{Username: c.Username, Token: fmt.Sprintf("%x", md5.Sum(t))}
	c.conn.Write(loginPacket.CreatePacket().ToBytesFast())
}
func (c *Client) packetReceiver() {
	for {
		header := make([]byte, 4)
		r, err := c.conn.Read(header)
		if err == io.EOF || r != 4 {
			fmt.Printf("Server closed the connection...\n")
			return
		}
		msgHeader := packets.ReadHeader(header)
		buff := make([]byte, msgHeader.Size)
		_, err = c.conn.Read(buff)
		if err != nil {
			fmt.Printf("Error while reading packet data %s", err.Error())
		}
		recvPacket, err := packets.FromBytes(msgHeader, buff)
		c.recv <- recvPacket
	}
}
func (c *Client) chat() string {
	cmd := c.readString()
	return cmd
}
func (c *Client) packetHandler() {
	for {
		select {
		case packet := <-c.recv:
			switch packet.Header.MsgType {
			case 2:
				msg, ok := packet.Data.(*packets.Message)
				//received message
				if ok {
					fmt.Printf("[%s] %s >>> %s: %s\n", formatTime(msg.Time), msg.Username, msg.Channel, msg.Message)
				}
			case 3:
				msg, ok := packet.Data.(*packets.SystemMessage)
				if ok {
					fmt.Printf("[%s] {{{SystemMessage}}}: %s}}}\n", formatTime(msg.Time), msg.Message)
				}
			}
		}
	}
}
func handleCommand(command string, args ...string) bool {
	switch command {
	case "channel":
		if len(args) == 0 {
			fmt.Printf("Active channel: %s\n", activeChannel)
		} else {
			activeChannel = args[0]
			fmt.Printf("Switched to channel: %s\n", activeChannel)
		}
		return true
	}
	return false
}
func main() {
	client, err := startClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client.login()
	go client.packetHandler()
	go client.packetReceiver()
	var command string
	for command != "/quit" {
		command = client.chat()
		var ok bool
		if len(command) > 0 && command[:1] == "/" {
			parsed := strings.Split(command, " ")
			command, args := parsed[0][1:], parsed[1:]
			ok = handleCommand(command, args...)
		}
		if !ok {
			msg := packets.Message{Username: client.Username, Message: command, Time: time.Now(), Channel: activeChannel}
			client.conn.Write(msg.CreatePacket().ToBytesFast())
		}
	}
	client.conn.Close()
}
