package bot

import "testing"

func TestCommand(t *testing.T) {
	cmd := NewCommand(
		"cmd <key>",
		"Description",
		func(conv Conversation) {

		},
	)

	if cmd.Command.Text != "cmd <key>" {
		t.Errorf("Command name does not match")
	}

	if cmd.Description != "Description" {
		t.Errorf("Command description does not match")
	}
}
