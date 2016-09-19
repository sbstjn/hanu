package bot

import (
	"testing"

	"github.com/sbstjn/platzhalter"
)

func TestConversation(t *testing.T) {
	command := platzhalter.NewCommand("cmd test <param>")

	message := Message{
		ID:   0,
		Text: "cmd test value",
	}

	conversation := NewConversation(&command, &message, nil)

	if conversation.Param("param") != "value" {
		t.Errorf("Failed to get correct value for param \"param\": %s != %s", conversation.Param("param"), "param")
	}

	conversation.Reply("example")
}
