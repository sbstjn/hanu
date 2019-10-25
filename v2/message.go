package hanu

import (
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

// MessageInterface defines the message interface
type MessageInterface interface {
	IsMessage() bool
	IsFrom(user string) bool
	IsHelpRequest() bool
	IsDirectMessage() bool
	IsMentionFor(user string) bool
	IsRelevantFor(user string) bool

	Text() string
	User() string
	Channel() string
}

func NewMessage(ev *slack.MessageEvent) Message {
	msg := Message{}
	msg.ChannelID = ev.Channel
	msg.Message = ev.Text
	msg.OriginalMessage = ev.Text
	msg.UserID = ev.Username
	msg.Type = ev.Type
	return msg
}

// Message is the Message structure for received and sent messages using Slack
type Message struct {
	ID              uint64
	Type            string
	ChannelID       string
	UserID          string
	Message         string
	OriginalMessage string
}

// Text returns the message text
func (m Message) Text() string {
	return m.Message
}

// Channel returns the channel ID
func (m Message) Channel() string {
	return m.ChannelID
}

// User returns the name of the user who sent the message
func (m Message) User() string {
	return m.UserID
}

// IsMessage checks if it is a Message or some other kind of processing information
func (m Message) IsMessage() bool {
	return true
}

// IsFrom checks the sender of the message
func (m Message) IsFrom(user string) bool {
	return m.UserID == user
}

// SetText updates the text of a message
func (m *Message) SetText(text string) {
	m.Message = text
}

// StripMention removes the mention from the message beginning
func (m *Message) StripMention(user string) string {
	prefix := "<@" + user + "> "
	text := m.Text()

	if strings.HasPrefix(text, prefix) {
		return text[len(prefix):len(text)]
	}

	return text
}

// StripLinkMarkup converts <http://google.com|google.com> into google.com etc.
// https://api.slack.com/docs/message-formatting#how_to_display_formatted_messages
func (m *Message) StripLinkMarkup() string {
	re := regexp.MustCompile("<(.*?)>")
	result := re.FindAllStringSubmatch(m.Text(), -1)
	text := m.Text()

	var link string
	for _, c := range result {
		link = c[len(c)-1]

		// Done change Channel, User or Specials tags
		if link[:2] == "#C" || link[:2] == "@U" || link[:1] == "!" {
			continue
		}

		url := link
		if strings.Contains(link, "|") {
			splits := strings.Split(link, "|")
			url = splits[1]
		}

		text = strings.Replace(text, "<"+link+">", url, -1)
	}

	return text
}

// IsHelpRequest checks if the user requests the help command
func (m Message) IsHelpRequest() bool {
	return strings.HasSuffix(m.Message, "help") || strings.HasPrefix(m.Message, "help")
}

// IsDirectMessage checks if the message is received using a direct messaging channel
func (m Message) IsDirectMessage() bool {
	return strings.HasPrefix(m.ChannelID, "D")
}

// IsMentionFor checks if the given user was mentioned with the message
func (m Message) IsMentionFor(user string) bool {
	return strings.HasPrefix(m.OriginalMessage, "<@"+user+">")
}

// IsRelevantFor checks if the message is relevant for a user
func (m Message) IsRelevantFor(user string) bool {
	return m.IsMessage() && !m.IsFrom(user) && (m.IsDirectMessage() || m.IsMentionFor(user))
}
