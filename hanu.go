package hanu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
		return nil, errors.New("Failed to connect to Slack RTM API")
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
		return nil, errors.New("Failed to connect to Websocket")
	}

	return b, nil
}

// Process incoming message
func (b *Bot) process(message *bot.SlackMessage) {
	if !message.IsRelevantFor(b.ID) {
		return
	}

	message.StripMention(b.ID)

	// Check if the message requests the auto-generated help command list
	// or if we need to search for a command matching the request
	if message.IsHelpRequest() {
		b.sendHelp(message)
	} else {
		b.searchCommand(message)
	}
}

// Search for a command matching the message
func (b *Bot) searchCommand(message *bot.SlackMessage) {
	for i := 0; i < len(b.Commands); i++ {
		if b.Commands[i].Command.Matches(message.Text) {
			b.Commands[i].Handler(bot.NewConversation(&b.Commands[i].Command, message, b.Socket))
		}
	}
}

// Send the response for a help request
func (b *Bot) sendHelp(message *bot.SlackMessage) {
	message.Text = "Thanks for asking! I can support you with those features:\n\n"

	for i := 0; i < len(b.Commands); i++ {
		message.Text = message.Text + "`" + b.Commands[i].Command.Text + "`"

		if b.Commands[i].Description != "" {
			message.Text = message.Text + " *–* " + b.Commands[i].Description
		}

		message.Text = message.Text + "\n"
	}

	if !message.IsDirectMessage() {
		message.Text = "<@" + message.User + ">: " + message.Text
	}

	websocket.JSON.Send(b.Socket, message)
}

// Listen for message on socket
func (b *Bot) Listen() {
	var msg bot.SlackMessage

	for {
		if websocket.JSON.Receive(b.Socket, &msg) != nil {
			log.Fatal("Error reading from Websocket")
		} else {
			b.process(&msg)
			msg = bot.SlackMessage{}
		}
	}
}

// Command adds a new command with custom handler
func (b *Bot) Command(command string, handler bot.CommandHandler) {
	b.Commands = append(b.Commands, bot.Command{
		Command:     platzhalter.NewCommand(command),
		Description: "",
		Handler:     handler,
	})
}

// Register registers a Command
func (b *Bot) Register(command bot.Command) {
	b.Commands = append(b.Commands, command)
}
