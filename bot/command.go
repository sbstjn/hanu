package bot

import (
	"github.com/sbstjn/platzhalter"
)

// CommandHandler is the interface for the handler function
type CommandHandler func(*Conversation)

// Command a command
type Command struct {
	Command     platzhalter.Command
	Description string
	Handler     CommandHandler
}

// NewCommand creates a new command
func NewCommand(command string, description string, handler CommandHandler) Command {
	return Command{
		Command:     platzhalter.NewCommand(command),
		Description: description,
		Handler:     handler,
	}
}
