package hanu

import (
	"fmt"

	"golang.org/x/net/websocket"

	"github.com/sbstjn/platzhalter"
)

// Conversation stores message, command and socket information and is passed
// to the handler function
type Conversation struct {
	message *Message
	command *platzhalter.Command
	socket  *websocket.Conn
}

// Reply sends message using the socket to Slack
func (c *Conversation) Reply(text string) {
	prefix := ""

	if !c.message.IsDirectMessage() {
		prefix = "<@" + c.message.User + ">: "
	}

	msg := c.message
	msg.Text = prefix + text
	fmt.Println("Reply: " + msg.Text)

	websocket.JSON.Send(c.socket, msg)
}

// Param gets a parameter value by name
func (c *Conversation) Param(name string) string {
	return c.command.Param(c.message.Text, name)
}
