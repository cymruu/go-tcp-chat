package server

import (
	"net"
)

type Client struct {
	Username string
	Token    string
	conn     net.Conn
}
