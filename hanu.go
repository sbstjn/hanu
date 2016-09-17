package hanu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sbstjn/hanu/bot"
	"github.com/sbstjn/platzhalter"

	"golang.org/x/net/websocket"
)

type handshakeResponseSelf struct {
	ID string `json:"id"`
}

type handshakeResponse struct {
	Ok    bool                  `json:"ok"`
	Error string                `json:"error"`
	URL   string                `json:"url"`
	Self  handshakeResponseSelf `json:"self"`
}

// Bot is the main object
type Bot struct {
	Socket   *websocket.Conn
	Token    string
	ID       string
	Commands []bot.Command
}

// New creates a new bot
func New(token string) (*Bot, error) {
	bot := Bot{
		Token: token,
	}

	return bot.Handshake()
}

// Handshake connects to the Slack API to get a socket connection
func (b *Bot) Handshake() (*Bot, error) {
	// Check for HTTP error on connection
	res, err := http.Get(fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", b.Token))
	if err != nil {
		return nil, errors.New("Failed to connect to slack")
	}

	// Check for HTTP status code
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed with HTTP Code: %d", res.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read body from response")
	}

	// Parse response
	var response handshakeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %s", body)
	}

	// Check for Slack error
	if !response.Ok {
		return nil, errors.New(response.Error)
	}

	// Assign Slack user ID
	b.ID = response.Self.ID

	// Connect to websocket
	b.Socket, err = websocket.Dial(response.URL, "", "https://api.slack.com/")
	if err != nil {
		return nil, errors.New("Failed to connect to socket")
	}

	return b, nil
}

// Process incoming message
func (b *Bot) process(msg *bot.Message) {
	if !msg.IsRelevantFor(b.ID) {
		return
	}

	msg.Text = strings.Trim(msg.Text, "<@"+b.ID+"> ")
	for i := 0; i < len(b.Commands); i++ {
		if b.Commands[i].Command.Matches(msg.Text) {
			b.Commands[i].Handler(bot.NewConversation(&b.Commands[i].Command, msg, b.Socket))
		}
	}
}

// Listen for message on socket
func (b *Bot) Listen() {
	var msg bot.Message

	for {
		if websocket.JSON.Receive(b.Socket, &msg) != nil {
			log.Fatal("Error reading from Websocket")
		} else {
			b.process(&msg)
			msg = bot.Message{}
		}
	}
}

// Command adds a new command with custom handler
func (b *Bot) Command(command string, handler bot.CommandHandler) {
	b.Commands = append(b.Commands, bot.Command{
		Command: platzhalter.NewCommand(command),
		Handler: handler,
	})
}

// Register registers a Command
func (b *Bot) Register(command bot.Command) {
	b.Commands = append(b.Commands, command)
}
