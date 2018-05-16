package main

import (
	"chat2/packets"
	"chat2/server"
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
	msg := &packets.Message{
		Username: "filip",
		Message:  "czesc",
		Time:     time.Now(),
	}
	p := msg.CreatePacket()
	fmt.Print(p.CreateHeader())
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
