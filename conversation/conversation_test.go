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
	command := allot.New("cmd test <param>")

	msg := message.Slack{
		ID: 0,
	}
	msg.SetText("cmd test value")

	match, _ := command.Match(msg.Text())

	conv := New(match, msg, nil)

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

	msg := message.Slack{
		ID: 0,
	}
	msg.SetText("cmd test value")

	match, _ := cmd.Match(msg.Text())

	conv := New(match, msg, &websocket.Conn{})
	conv.SetConnection(ConnectionMock{})

	conv.Reply("example")
}
