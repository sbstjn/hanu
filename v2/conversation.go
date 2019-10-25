package hanu

import (
	"github.com/sbstjn/allot"
)

// Convo is a shorthand for ConversationInterface
type Convo ConversationInterface

// ConversationInterface is the interface for a conversation
type ConversationInterface interface {
	Integer(name string) (int, error)
	String(name string) (string, error)
	Reply(text string, a ...interface{})
	Match(position int) (string, error)
	Message() MessageInterface
}

// Sayer is an object that can talk in the channel
type Sayer interface {
	Say(string, string, ...interface{})
}

// Conversation stores message, command and socket information and is passed
// to the handler function
type Conversation struct {
	message Message
	match   allot.MatchInterface
	bot     Sayer
}

// Message returns the convos message
func (c *Conversation) Message() MessageInterface {
	return c.message
}

// Reply sends message using the socket to Slack
func (c *Conversation) Reply(text string, a ...interface{}) {
	prefix := ""

	if !c.message.IsDirectMessage() {
		prefix = "<@" + c.message.User() + ">: "
	}

	c.bot.Say(c.Message().Channel(), prefix+text, a...)
}

// String return string paramter
func (c Conversation) String(name string) (string, error) {
	return c.match.String(name)
}

// Integer returns integer parameter
func (c Conversation) Integer(name string) (int, error) {
	return c.match.Integer(name)
}

// Match returns the parameter at the position
func (c Conversation) Match(position int) (string, error) {
	return c.match.Match(position)
}

// NewConversation returns a Conversation struct
func NewConversation(match allot.MatchInterface, msg Message, bot Sayer) ConversationInterface {
	conv := &Conversation{
		message: msg,
		match:   match,
		bot:     bot,
	}

	return conv
}
