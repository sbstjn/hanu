package conversation

import (
	"testing"

	"github.com/sbstjn/hanu/message"
	"github.com/sbstjn/platzhalter"
)

func TestConversation(t *testing.T) {
	command := platzhalter.NewCommand("cmd test <param>")

	msg := message.Slack{
		ID: 0,
	}
	msg.SetText("cmd test value")

	conv := New(&command, &msg, nil)

	if conv.Param("param") != "value" {
		t.Errorf("Failed to get correct value for param \"param\": %s != %s", conv.Param("param"), "param")
	}

	conv.Reply("example")
}
