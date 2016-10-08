package message

import "testing"

func TestMessage(t *testing.T) {
	msg := Slack{
		ID:   0,
		Type: "message",
	}

	if !msg.IsMessage() {
		t.Errorf("IsMessage() must be true")
	}

	if msg.IsDirectMessage() {
		t.Errorf("msg.IsDirectMessage() must be false")
	}

	if msg.IsMentionFor("") {
		t.Errorf("msg.IsMentionFor() must be false")
	}

	if msg.IsRelevantFor("user") {
		t.Errorf("msg.IsRelevantFor() must be true")
	}
}

func TestHelpMessage(t *testing.T) {
	msg := Slack{}
	msg.SetText("help")

	if !msg.IsHelpRequest() {
		t.Errorf("msg.IsHelpRequest() must be true")
	}
}

func StripMention(t *testing.T) {
	msg := Slack{}
	msg.SetText("<@test> help")

	msg.StripMention("test")

	if msg.Text() != "help" {
		t.Errorf("msg.Text must be 'help'")
	}
}
