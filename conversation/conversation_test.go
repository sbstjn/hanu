package conversation

import (
	"testing"

	"golang.org/x/net/websocket"

	"github.com/sbstjn/allot"
	"github.com/sbstjn/hanu/message"
)

type ConnectionMock struct{}

func (c ConnectionMock) Send(ws *websocket.Conn, v interface{}) (err error) {
	return nil
}

func TestConversation(t *testing.T) {
	command := allot.NewCommand("cmd test <param>")

	msg := message.Slack{
		ID: 0,
	}
	msg.SetText("cmd test value")

	conv := New(&command, msg, nil)

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
	cmd := allot.NewCommand("cmd test <param>")

	msg := message.Slack{
		ID: 0,
	}
	msg.SetText("cmd test value")

	conv := New(&cmd, msg, &websocket.Conn{})
	conv.SetConnection(ConnectionMock{})

	conv.Reply("example")
}
