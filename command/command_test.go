package command

import (
	"testing"

	"github.com/sbstjn/hanu/conversation"
	"github.com/sbstjn/hanu/message"
)

func TestCommand(t *testing.T) {
	cmd := New(
		"cmd <key>",
		"Description",
		func(conv conversation.Interface) {

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
	cmd := New(
		"cmd <key>",
		"Description",
		func(conv conversation.Interface) {
			str, _ := conv.String("key")
			if str != "name" {
				t.Errorf("param <key> should have value \"name\"")
			}
		},
	)

	msg := message.Slack{}
	msg.SetText("cmd name")

	conv := conversation.New(cmd.Get(), msg, nil)
	cmd.Handle(conv)
}
