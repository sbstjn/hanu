package bot

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
func (c *Conversation) Reply(text string, a ...interface{}) {
	prefix := ""

	if !c.message.IsDirectMessage() {
		prefix = "<@" + c.message.User + ">: "
	}

	msg := c.message
	msg.Text = prefix + fmt.Sprintf(text, a...)

	if c.socket != nil {
		websocket.JSON.Send(c.socket, msg)
	}
}

// Param gets a parameter value by name
func (c *Conversation) Param(name string) string {
	return c.command.Param(c.message.Text, name)
}

// NewConversation returns a Conversation struct
func NewConversation(command *platzhalter.Command, message *Message, socket *websocket.Conn) *Conversation {
	return &Conversation{
		message: message,
		command: command,
		socket:  socket,
	}
}
