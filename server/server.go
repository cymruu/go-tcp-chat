package server

import (
	"bufio"
	"chat2/packets"
	"encoding/gob"
	"fmt"
	"net"
)

type Server struct {
	listener     net.Listener
	connnections []*net.Conn
	incoming     []packets.IPacket
}

func (s *Server) Listen(addr string) error {
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("Server started [%s]... waiting for connections\n", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		s.connnections = append(s.connnections, &conn)
		go s.handleConnection(&conn)
	}
}
func (s *Server) handleConnection(conn *net.Conn) {
	fmt.Printf("Hadling new connection %p\n", conn)
	fmt.Printf("connections %v\n", s.connnections)
	rw := bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))
	dec := gob.NewDecoder(rw)
	for {
		header := make([]byte, 4)
		_, err := rw.Read(header)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		msgSize, msgType := packets.ReadHeader(header)
		fmt.Printf("Received packet length: %v Type: %d\n", msgSize, msgType)
		// data := make([]byte, dataLength)
		// rw.Read(data)
		// fmt.Printf("received packet... %v", data)
		var packet packets.Message
		err = dec.Decode(&packet)
		if err != nil {
			fmt.Printf("Error while decoding packet %s", err.Error())
		}
		fmt.Println(packet)
	}
}
func CreateServer() *Server {
	return &Server{}
}
