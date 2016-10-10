package command

import (
	"github.com/sbstjn/hanu/conversation"
	"github.com/sbstjn/platzhalter"
)

// Handler is the interface for the handler function
type Handler func(conversation.Interface)

// Interface defines a command interface
type Interface interface {
	Get() *platzhalter.Command
	Set(cmd platzhalter.Command)
	Description() string
	SetDescription(text string)
	Handle(conv conversation.Interface)
	SetHandler(handler Handler)
}

// Command a command
type Command struct {
	command     *platzhalter.Command
	description string
	handler     Handler
}

// SetHandler sets the handler
func (c *Command) SetHandler(handler Handler) {
	c.handler = handler
}

// Description returns the description
func (c *Command) Description() string {
	return c.description
}

// SetDescription sets the description
func (c *Command) SetDescription(text string) {
	c.description = text
}

// Handle calls the command's handler
func (c *Command) Handle(conv conversation.Interface) {
	go c.handler(conv)
}

// Get returns the platzhalter command
func (c *Command) Get() *platzhalter.Command {
	return c.command
}

// Set defines the platzhalter command
func (c *Command) Set(cmd platzhalter.Command) {
	c.command = &cmd
}

// New creates a new command
func New(command string, description string, handler Handler) Interface {
	cmd := &Command{}
	cmd.Set(platzhalter.NewCommand(command))
	cmd.SetDescription(description)
	cmd.SetHandler(handler)

	return cmd
}
