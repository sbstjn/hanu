package hanu

import "testing"

func TestMessage(t *testing.T) {
	msg := Message{
		ID:   0,
		Type: "message",
	}

	if !msg.IsMessage() {
		t.Errorf("IsMessage() must be true")
	}

	if msg.IsDirectMessage() {
		t.Errorf("msg.IsDirectMessage() must be false")
	}
}
