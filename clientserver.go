package main

import (
	"bufio"
	"chat2/packets"
	"chat2/server"
	"crypto/md5"
	"fmt"
	"io/ioutil"
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
	reader := bufio.NewReader(client)
	data, _ := ioutil.ReadAll(reader)
	fmt.Print(data)
	client.Close()
}
func startServer() {
	srv := server.CreateServer()
	srv.Listen(":3300")
}
func main() {
	go startServer()
	go client()
	srv := server.CreateServer()
	srv.Listen(":3300")
}
