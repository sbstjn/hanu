package hanu

import (
	"github.com/sbstjn/allot"
)

// Handler is the interface for the handler function
type Handler func(ConversationInterface)

// CommandInterface defines a command interface
type CommandInterface interface {
	Get() allot.CommandInterface
	Description() string
	Handle(conv ConversationInterface)
}

// Command a command
type Command struct {
	command     allot.CommandInterface
	description string
	handler     Handler
}

// SetHandler sets the handler
func (c *Command) SetHandler(handler Handler) {
	c.handler = handler
}

// Description returns the description
func (c Command) Description() string {
	return c.description
}

// SetDescription sets the description
func (c *Command) SetDescription(text string) {
	c.description = text
}

// Handle calls the command's handler
func (c Command) Handle(conv ConversationInterface) {
	go c.handler(conv)
}

// Get returns the platzhalter command
func (c Command) Get() allot.CommandInterface {
	return c.command
}

// Set defines the platzhalter command
func (c *Command) Set(cmd allot.CommandInterface) {
	c.command = cmd
}

// NewCommand creates a new command
func NewCommand(text string, description string, handler Handler) Command {
	cmd := Command{}
	cmd.Set(allot.New(text))
	cmd.SetDescription(description)
	cmd.SetHandler(handler)

	return cmd
}
