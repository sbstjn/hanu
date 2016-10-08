package conversation

import (
	"fmt"

	"golang.org/x/net/websocket"

	"github.com/sbstjn/hanu/message"
	"github.com/sbstjn/platzhalter"
)

// Interface is the interface for a conversation
type Interface interface {
	Reply(text string, a ...interface{})
	Param(name string) string
	SetConnection(connection Connection)

	send(msg message.Interface)
}

// Connection is the needed interface for a connection
type Connection interface {
	Send(ws *websocket.Conn, v interface{}) (err error)
}

// Slack stores message, command and socket information and is passed
// to the handler function
type Slack struct {
	message message.Interface
	command *platzhalter.Command
	socket  *websocket.Conn

	connection Connection
}

func (c *Slack) send(msg message.Interface) {
	if c.socket != nil {
		c.connection.Send(c.socket, msg)
	}
}

// SetConnection sets the conversation connection
func (c *Slack) SetConnection(connection Connection) {
	c.connection = connection
}

// Reply sends message using the socket to Slack
func (c *Slack) Reply(text string, a ...interface{}) {
	prefix := ""

	if !c.message.IsDirectMessage() {
		prefix = "<@" + c.message.User() + ">: "
	}

	msg := c.message
	msg.SetText(prefix + fmt.Sprintf(text, a...))

	c.send(msg)
}

// Param gets a parameter value by name
func (c *Slack) Param(name string) string {
	return c.command.Param(c.message.Text(), name)
}

// New returns a Conversation struct
func New(command *platzhalter.Command, msg message.Interface, socket *websocket.Conn) Interface {
	conv := &Slack{
		message: msg,
		command: command,
		socket:  socket,
	}

	conv.SetConnection(websocket.JSON)

	return conv
}
