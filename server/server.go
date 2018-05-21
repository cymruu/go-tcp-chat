package server

import (
	"bufio"
	"chat2/packets"
	"fmt"
	"io"
	"net"
)

type handleFunction func(*Server, *SocketMessage)

var handlers = map[uint16]handleFunction{
	1: AuthorizationFunction,
	2: MessageFunction,
}

func MessageFunction(s *Server, sm *SocketMessage) {
	msg, ok := sm.msg.Data.(*packets.Message)
	if ok != true {
		fmt.Printf("Cant assert Message packet\n")
		return
	}
	fmt.Printf("Received message %v\n", msg)
	for _, c := range s.connnections {
		fmt.Printf("Sending msg to %s\n", c.Username)
		c.conn.Write(msg.CreatePacket().ToBytesFast())
	}
}
func AuthorizationFunction(s *Server, sm *SocketMessage) {
	msg, ok := sm.msg.Data.(*packets.Authorization)
	if ok != true {
		fmt.Printf("Cant assert Authorization packet")
		return
	}
	fmt.Printf("Authorizing socket as %s\n", msg.Username)
	sm.c.Username = msg.Username
	sm.c.Token = msg.Token
	return
}

type Server struct {
	listener     net.Listener
	connnections []*Client
	incoming     chan SocketMessage
	quitChannel  chan bool
}
type SocketMessage struct {
	c   *Client
	msg packets.Packet
}

func (s *Server) Listen(addr string) error {
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("Server started [%s]... waiting for connections\n", s.listener.Addr())
	go s.readMessages()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		go s.handleConnection(conn)
	}
}
func (s *Server) readMessages() {
	for {
		select {
		case sm := <-s.incoming:
			fmt.Printf("Channel handler received message %+v \n", sm)
			handler, ok := handlers[sm.msg.Header.MsgType]

			if ok == true {
				go handler(s, &sm)
			}
		}
	}
}
func (s *Server) handleConnection(conn net.Conn) {
	fmt.Printf("Hadling new connection %p\n", conn)
	fmt.Printf("connections %v\n", s.connnections)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client := &Client{conn: conn}
	s.connnections = append(s.connnections, client)
	for {
		header := make([]byte, 4)
		r, err := rw.Read(header)
		if err == io.EOF || r != 4 {
			return
		}
		msgHeader := packets.ReadHeader(header)
		buff := make([]byte, msgHeader.Size)
		_, err = rw.Read(buff)
		if err != nil {
			fmt.Printf("Error while reading packet data %s", err.Error())
		}
		recvPacket, err := packets.FromBytes(msgHeader, buff)
		sm := SocketMessage{c: client, msg: recvPacket}
		s.incoming <- sm
	}
}
func CreateServer() *Server {
	return &Server{
		incoming:    make(chan SocketMessage),
		quitChannel: make(chan bool, 1),
	}
}
