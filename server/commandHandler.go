package server

import (
	"chat2/packets"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type commandHandler func(*Server, *Client, ...string)

func handleQuit(s *Server, c *Client, args ...string) {
	for i, cli := range s.connnections {
		if c == cli {
			s.connnections = append(s.connnections[:i], s.connnections[i+1:]...)
			s.broadcast(&packets.SystemMessage{Message: fmt.Sprintf("%s has left", c.Username)})
			c.conn.Close()
		}
	}
}

func handleRoll(s *Server, c *Client, argsD ...string) {
	args := []string{"6", "1"} //roll default dice one time
	copy(args, argsD)
	sides, err := strconv.Atoi(args[0])
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		c.sendData(&packets.SystemMessage{Message: "Bad arguments... usage /roll diceSides amountOfRolls (max 10)"})
		return
	}
	if amount > 10 {
		amount = 10
	}
	fmt.Printf("rolling %d times sides %d", amount, sides)
	for i := 0; i < amount; i++ {
		roll := 1 + rand.Intn(sides)
		fmt.Println(roll)
		s.broadcast(&packets.SystemMessage{Message: fmt.Sprintf("%s uses the dice and rolls %d", c.Username, roll)})
	}

}

var commandHandlers = map[string]commandHandler{
	"quit": handleQuit,
	"roll": handleRoll,
}

func (s *Server) handleCommand(c *Client, msg *packets.Message) {
	parsed := strings.Split(msg.Message, " ")
	command, args := parsed[0][1:], parsed[1:]
	fmt.Printf("COMMAND %s with %q\n", command, args)
	handler, ok := commandHandlers[command]
	if ok {
		handler(s, c, args...)
	} else {
		c.sendData(&packets.SystemMessage{Message: fmt.Sprintf("%s command doesn't exists", command)})
	}
}
