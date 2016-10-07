package bot

import "testing"

func TestMessage(t *testing.T) {
	msg := SlackMessage{
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
