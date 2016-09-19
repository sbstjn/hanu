package bot

import "strings"

// Message is the Message structure for received and sent messages
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

// IsMessage checks if it is a Message or some other kind of processing information
func (m *Message) IsMessage() bool {
	return m.Type == "message"
}

// IsFrom checks the sender of the message
func (m *Message) IsFrom(user string) bool {
	return m.User == user
}

// IsDirectMessage checks if the message is received using a direct messaging channel
func (m *Message) IsDirectMessage() bool {
	return strings.HasPrefix(m.Channel, "D")
}

// IsMentionFor checks if the given user was mentioned with the message
func (m *Message) IsMentionFor(user string) bool {
	return strings.HasPrefix(m.Text, "<@"+user+">")
}

// IsRelevantFor checks if the message is relevant for a user
func (m *Message) IsRelevantFor(user string) bool {
	return m.IsMessage() && !m.IsFrom(user) && (m.IsDirectMessage() || m.IsMentionFor(user))
}
