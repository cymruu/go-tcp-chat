package main

import (
	"chat2/server"
)

func main() {
	srv := server.CreateServer()
	srv.Listen(":3300")
}
