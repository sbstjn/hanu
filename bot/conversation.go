package bot

import (
	"fmt"

	"golang.org/x/net/websocket"

	"github.com/sbstjn/platzhalter"
)

// Conversation is the interface for a conversation
type Conversation interface {
	Reply(text string, a ...interface{})
	Param(name string) string

	send(msg *SlackMessage)
}

// SlackConversation stores message, command and socket information and is passed
// to the handler function
type SlackConversation struct {
	message *SlackMessage
	command *platzhalter.Command
	socket  *websocket.Conn
}

func (c *SlackConversation) send(msg *SlackMessage) {
	if c.socket != nil {
		websocket.JSON.Send(c.socket, msg)
	}
}

// Reply sends message using the socket to Slack
func (c *SlackConversation) Reply(text string, a ...interface{}) {
	prefix := ""

	if !c.message.IsDirectMessage() {
		prefix = "<@" + c.message.User + ">: "
	}

	msg := c.message
	msg.Text = prefix + fmt.Sprintf(text, a...)

	c.send(msg)
}

// Param gets a parameter value by name
func (c *SlackConversation) Param(name string) string {
	return c.command.Param(c.message.Text, name)
}

// NewConversation returns a Conversation struct
func NewConversation(command *platzhalter.Command, message *SlackMessage, socket *websocket.Conn) *SlackConversation {
	return &SlackConversation{
		message: message,
		command: command,
		socket:  socket,
	}
}
