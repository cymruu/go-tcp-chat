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
	if msg.Message == "/quit" {
		for i, c := range s.connnections {
			if c == sm.c {
				s.connnections = append(s.connnections[:i], s.connnections[i+1:]...)
				c.conn.Close()
				leftMessage := packets.SystemMessage{Message: fmt.Sprintf("%s has left", c.Username)}
				s.broadcast(&leftMessage)
			}
		}
		return
	}
	s.broadcast(msg)
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
	welcomeMessage := packets.SystemMessage{Message: fmt.Sprintf("%p is now known as %s", sm.c, sm.c.Username)}
	s.broadcast(&welcomeMessage)
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
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client := &Client{conn: conn}
	s.connnections = append(s.connnections, client)
	for {
		header := make([]byte, 4)
		_, err := rw.Read(header)
		if err == io.EOF {
			return
		}
		msgHeader := packets.ReadHeader(header)
		buff := make([]byte, msgHeader.Size)
		_, err = rw.Read(buff)
		if err != nil {
			fmt.Printf("Error while reading packet data %s", err.Error())
			conn.Close()
		}
		recvPacket, err := packets.FromBytes(msgHeader, buff)
		if err == nil {
			sm := SocketMessage{c: client, msg: recvPacket}
			s.incoming <- sm
		}
	}
}
func (s *Server) broadcast(packet packets.IPacketData) {
	for _, c := range s.connnections {
		if len(c.Username) > 0 { //send only to authorized clients
			fmt.Printf("Sending msg %v to %s\n", packet, c.Username)
			c.sendData(packet)
		}
	}
}
func CreateServer() *Server {
	return &Server{
		incoming:    make(chan SocketMessage),
		quitChannel: make(chan bool, 1),
	}
}
