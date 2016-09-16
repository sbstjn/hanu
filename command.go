package hanu

import "github.com/sbstjn/platzhalter"

// CommandHandler is the interface for the handler function
type CommandHandler func(*Conversation)

// Command a command
type Command struct {
	Command platzhalter.Command
	Handler CommandHandler
}
