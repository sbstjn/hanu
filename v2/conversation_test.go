package hanu

import (
	"testing"

	"github.com/sbstjn/allot"
)

type SayerMock struct {
	msg  string
	ch   string
	args []interface{}
}

func (sm SayerMock) Say(ch, msg string, a ...interface{}) {
	sm.ch = ch
	sm.msg = msg
	sm.args = a
}

func TestConversation(t *testing.T) {
	command := allot.New("cmd test <param>")

	msg := Message{
		ID: 0,
	}
	msg.SetText("cmd test value")

	match, _ := command.Match(msg.Text())

	conv := NewConversation(match, msg, &SayerMock{})

	str, err := conv.String("param")

	if err != nil {
		t.Errorf("Failed to get correct value for param \"param\"")
	}

	if str != "value" {
		t.Errorf("Failed to get correct value for param \"param\": %s != %s", str, "value")
	}

	conv.Reply("example")
}

func TestConnect(t *testing.T) {
	cmd := allot.New("cmd test <param>")

	msg := Message{
		ID: 0,
	}
	msg.SetText("cmd test value")

	match, _ := cmd.Match(msg.Text())

	conv := NewConversation(match, msg, &SayerMock{})

	conv.Reply("example")
}
