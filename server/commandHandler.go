package server

import (
	"chat2/packets"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type commandHandler func(*Server, *Channel, *Client, ...string)

func handleQuit(s *Server, channel *Channel, c *Client, args ...string) {
	for i, cli := range s.connnections {
		if c == cli {
			s.connnections = append(s.connnections[:i], s.connnections[i+1:]...)
			s.broadcast(&packets.SystemMessage{Message: fmt.Sprintf("%s has disconnected", c.Username)})
			c.conn.Close()
		}
	}
}

func handleRoll(s *Server, channel *Channel, c *Client, argsD ...string) {
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
	fmt.Printf("rolling %d times sides %d\n", amount, sides)
	for i := 0; i < amount; i++ {
		roll := 1 + rand.Intn(sides)
		channel.broadcast(&packets.SystemMessage{Message: fmt.Sprintf("%s uses the dice and rolls %d", c.Username, roll), Channel: channel.name})
	}
}
func handleJoin(s *Server, _ *Channel, c *Client, argsD ...string) {
	args := []string{""}
	copy(args, argsD)
	channelName := args[0]
	if channelName == "" {
		c.sendData(&packets.SystemMessage{Message: fmt.Sprintf("Usage: /join <channelName>")})
		return
	}
	channel, ok := s.channels[channelName]
	if !ok {
		channel = s.CreateChannel(channelName)
	}
	channel.join(c)
}

func handleLeave(s *Server, _ *Channel, c *Client, argsD ...string) {
	args := []string{""}
	copy(args, argsD)
	channelName := args[0]
	channel, ok := s.channels[channelName]
	if ok {
		channel.leave(c)
	}
}
func handleOnline(s *Server, channel *Channel, c *Client, _ ...string) {
	online := fmt.Sprintf("List of users in room: %s: ", channel.name)
	for _, participant := range channel.participants {
		online += "@" + participant.Username + "," + " "
	}
	c.sendData(&packets.SystemMessage{Message: online})
}

var commandHandlers = map[string]commandHandler{
	"quit":   handleQuit,
	"roll":   handleRoll,
	"join":   handleJoin,
	"leave":  handleLeave,
	"online": handleOnline,
}

func (s *Server) handleCommand(c *Client, msg *packets.Message) {
	parsed := strings.Split(msg.Message, " ")
	command, args := parsed[0][1:], parsed[1:]
	fmt.Printf("COMMAND %s args: %q\n", command, args)
	channel, _ := s.channels[msg.Channel]
	handler, ok := commandHandlers[command]
	if ok {
		handler(s, channel, c, args...)
	} else {
		c.sendData(&packets.SystemMessage{Message: fmt.Sprintf("%s command doesn't exists", command)})
	}
}
