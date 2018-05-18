package main

import (
	"chat2/packets"
	"chat2/server"
	"crypto/md5"
	"fmt"
	"net"
	"time"
)

func client() {
	time.Sleep(time.Second)
	client, err := net.Dial("tcp", ":3300")
	if err != nil {
		fmt.Print(err.Error())
	}
	t, _ := time.Now().MarshalBinary()
	auth := &packets.Authorization{
		Username: "filipek",
		Token:    fmt.Sprintf("%x", md5.Sum(t)),
	}
	msg := &packets.Message{
		Username: "filip",
		Message:  "czesc",
		Time:     time.Now(),
	}
	p := msg.CreatePacket()
	client.Write(auth.CreatePacket().ToBytes())
	client.Write(p.ToBytesFast())
	time.Sleep(time.Second)
	client.Write(p.ToBytesFast())
	client.Close()
}
func main() {
	go client()
	srv := server.CreateServer()
	srv.Listen(":3300")
}
