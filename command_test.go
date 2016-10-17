package hanu

import (
	"testing"
)

func TestCommand(t *testing.T) {
	cmd := NewCommand(
		"cmd <key>",
		"Description",
		func(conv ConversationInterface) {

		},
	)

	if cmd.Get().Text() != "cmd <key>" {
		t.Errorf("Command name does not match")
	}

	if cmd.Description() != "Description" {
		t.Errorf("Command description does not match")
	}
}

func TestHandle(t *testing.T) {
	cmd := NewCommand(
		"cmd <key>",
		"Description",
		func(conv ConversationInterface) {
			str, _ := conv.String("key")
			if str != "name" {
				t.Errorf("param <key> should have value \"name\"")
			}
		},
	)

	msg := Message{}
	msg.SetText("cmd name")

	match, _ := cmd.Get().Match(msg.Text())

	conv := NewConversation(match, msg, nil)
	cmd.Handle(conv)
}
