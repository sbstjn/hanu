package hanu

import (
	"regexp"
	"strings"

	"github.com/slack-go/slack/slackevents"
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
	UserID() string
	Channel() string
}

// NewMessage will create a new message object when given
// a slack message object
func NewMessage(ev interface{}) Message {
	m := Message{}
	sm := new(slackMessage)
	switch v := ev.(type) {
	case *slackevents.AppMentionEvent:
		sm.AppMentionEvent = v
	case *slackevents.MessageEvent:
		sm.MessageEvent = v
	}
	m.Message = sm.Text()
	m.ChannelID = sm.Channel()
	m.Username = sm.User()
	m.userID = sm.User() // why?
	return m
}

type slackMessage struct {
	*slackevents.MessageEvent
	*slackevents.AppMentionEvent
}

func (m slackMessage) Text() string {
	if m.MessageEvent != nil {
		return m.MessageEvent.Text
	}
	if m.AppMentionEvent != nil {
		return m.AppMentionEvent.Text
	}
	return ""
}

func (m slackMessage) User() string {
	if m.MessageEvent != nil {
		return m.MessageEvent.User
	}
	if m.AppMentionEvent != nil {
		return m.AppMentionEvent.User
	}
	return ""
}

func (m slackMessage) Channel() string {
	if m.MessageEvent != nil {
		return m.MessageEvent.Channel
	}
	if m.AppMentionEvent != nil {
		return m.AppMentionEvent.Channel
	}
	return ""
}

// Message is the Message structure for received and
// sent messages using Slack
type Message struct {
	*slackMessage
	Message   string
	ChannelID string
	userID    string
}

// Text returns the message text
func (m Message) Text() string {
	return m.Message
}

// Channel returns the channel ID
func (m Message) Channel() string {
	return m.ChannelID
}

func (m Message) UserID() string { //why?
	return m.userID
}

// User returns the name of the user who sent the message
func (m Message) User() string {
	return m.Username
}

// IsMessage checks if it is a Message or some other kind of processing information
func (m Message) IsMessage() bool {
	return true
}

// IsFrom checks the sender of the message
func (m Message) IsFrom(user string) bool {
	return m.User() == user
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
		m.Message = text[len(prefix):]
	}

	return m.Text()
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

	m.Message = text
	return text
}

// IsHelpRequest checks if the user requests the help command
func (m Message) IsHelpRequest() bool {
	return strings.HasSuffix(m.Message, "help") || strings.HasPrefix(m.Message, "help")
}

// IsDirectMessage checks if the message is received using a direct messaging channel
func (m Message) IsDirectMessage() bool {
	return strings.HasPrefix(m.Channel(), "D")
}

// IsMentionFor checks if the given user was mentioned with the message
func (m Message) IsMentionFor(user string) bool {
	return strings.HasPrefix(m.MessageEvent.Text, "<@"+user+">")
}

// IsRelevantFor checks if the message is relevant for a user
func (m Message) IsRelevantFor(user string) bool {
	return m.IsMessage() && !m.IsFrom(user) && (m.IsDirectMessage() || m.IsMentionFor(user))
}
