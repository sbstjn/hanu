package command

import (
	"testing"

	"github.com/sbstjn/hanu/conversation"
)

func TestCommand(t *testing.T) {
	cmd := New(
		"cmd <key>",
		"Description",
		func(conv conversation.Interface) {

		},
	)

	if cmd.Get().Text != "cmd <key>" {
		t.Errorf("Command name does not match")
	}

	if cmd.Description() != "Description" {
		t.Errorf("Command description does not match")
	}
}
