package message

import "strings"

// Interface defines the message interface
type Interface interface {
	IsMessage() bool
	IsFrom(user string) bool
	IsHelpRequest() bool
	IsDirectMessage() bool
	IsMentionFor(user string) bool
	IsRelevantFor(user string) bool

	Text() string
	User() string
}

// Slack is the Message structure for received and sent messages using Slack
type Slack struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	UserID  string `json:"user"`
	Message string `json:"text"`
}

// Text returns the message text
func (m Slack) Text() string {
	return m.Message
}

// User returns the message text
func (m Slack) User() string {
	return m.UserID
}

// IsMessage checks if it is a Message or some other kind of processing information
func (m Slack) IsMessage() bool {
	return m.Type == "message"
}

// IsFrom checks the sender of the message
func (m Slack) IsFrom(user string) bool {
	return m.UserID == user
}

// SetText updates the text of a message
func (m *Slack) SetText(text string) {
	m.Message = text
}

// StripMention removes the mention from the message beginning
func (m *Slack) StripMention(user string) {
	m.SetText(strings.Trim(m.Message, "<@"+user+"> "))
}

// IsHelpRequest checks if the user requests the help command
func (m Slack) IsHelpRequest() bool {
	return strings.HasSuffix(m.Message, "help") || strings.HasPrefix(m.Message, "help")
}

// IsDirectMessage checks if the message is received using a direct messaging channel
func (m Slack) IsDirectMessage() bool {
	return strings.HasPrefix(m.Channel, "D")
}

// IsMentionFor checks if the given user was mentioned with the message
func (m Slack) IsMentionFor(user string) bool {
	return strings.HasPrefix(m.Message, "<@"+user+">")
}

// IsRelevantFor checks if the message is relevant for a user
func (m Slack) IsRelevantFor(user string) bool {
	return m.IsMessage() && !m.IsFrom(user) && (m.IsDirectMessage() || m.IsMentionFor(user))
}
