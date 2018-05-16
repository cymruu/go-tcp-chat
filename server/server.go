package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

type Server struct {
	listener     net.Listener
	connnections []*net.Conn
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
	fmt.Printf("connections %v", s.connnections)
	rw := bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))
	for {
		dataLength := make([]byte, 2)
		_, err := rw.Read(dataLength)
		if err != nil {
			continue
		}
		u16 := binary.BigEndian.Uint16(dataLength)
		fmt.Printf("Received packet length: %v", u16)
		data := make([]byte, u16)
		rw.Read(data)
		fmt.Printf("received packet... %v", data)
	}
}
func CreateServer() *Server {
	return &Server{}
}
