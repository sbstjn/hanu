package hanu

import (
	"context"
	"fmt"

	"github.com/nlopes/slack"
)

// Bot is the main object
type Bot struct {
	RTM               *slack.RTM
	ID                string
	Commands          []CommandInterface
	ReplyOnly         bool
	CmdPrefix         string
	unknownCmdHandler Handler
}

// New creates a new bot
func New(token string) (*Bot, error) {
	api := slack.New(token)
	id, err := api.GetUserIdentity()
	if err != nil {
		return nil, err
	}

	rtm := api.NewRTM()
	bot := &Bot{RTM: rtm, ID: id.User.ID}
	return bot, nil
}

// SetCommandPrefix will set thing that must be prefixed to the command,
// there is no prefix by default but one could set it to "!" for instance
func (b *Bot) SetCommandPrefix(pfx string) *Bot {
	b.CmdPrefix = pfx
	return b
}

// SetReplyOnly will make the bot only respond to messages it is mentioned in
func (b *Bot) SetReplyOnly(ro bool) *Bot {
	b.ReplyOnly = ro
	return b
}

// Process incoming message
func (b *Bot) process(msg Message) {
	// Strip @BotName from public message
	msg.SetText(msg.StripMention(b.ID))
	// Strip Slack's link markup
	msg.SetText(msg.StripLinkMarkup())

	// Only send auto-generated help command list if directly mentioned
	if msg.IsRelevantFor(b.ID) && msg.IsHelpRequest() {
		b.sendHelp(msg)
		return
	}

	// if bot can only reply, ensure we were mentioned
	if b.ReplyOnly && !msg.IsRelevantFor(b.ID) {
		return
	}

	handled := b.searchCommand(msg)
	if !handled && b.ReplyOnly {
		if b.unknownCmdHandler != nil {
			b.unknownCmdHandler(NewConversation(dummyMatch{}, msg, b))
		}
	}
}

// Search for a command matching the message
func (b *Bot) searchCommand(msg Message) bool {
	var cmd CommandInterface

	for i := 0; i < len(b.Commands); i++ {
		cmd = b.Commands[i]

		match, err := cmd.Get().Match(msg.Text())
		if err == nil {
			cmd.Handle(NewConversation(match, msg, b))
			return true
		}
	}

	return false
}

// Channel will return a channel that the bot can talk in
func (b *Bot) Channel(id string) Channel {
	return Channel{b, id}
}

// Say will cause the bot to say something in the specified channel
func (b *Bot) Say(channel, msg string, a ...interface{}) {
	b.send(Message{ChannelID: channel, Message: fmt.Sprintf(msg, a...)})
}

func (b *Bot) send(msg MessageInterface) {
	b.RTM.SendMessage(&slack.OutgoingMessage{
		Channel: msg.Channel(),
		Text:    msg.Text(),
		Type:    "message",
	})
}

// BuildHelpText will build the help text
func (b *Bot) BuildHelpText() string {
	var cmd CommandInterface
	help := "The available commands are:\n\n"

	for i := 0; i < len(b.Commands); i++ {
		cmd = b.Commands[i]

		help = help + "`" + b.CmdPrefix + cmd.Get().Text() + "`"
		if cmd.Description() != "" {
			help = help + " *–* " + cmd.Description()
		}

		help = help + "\n"
	}

	return help
}

// sendHelp will send help to the channel and user in the given message
func (b *Bot) sendHelp(msg MessageInterface) {
	help := b.BuildHelpText()

	if !msg.IsDirectMessage() {
		help = "<@" + msg.User() + ">: " + help
	}

	b.Say(msg.Channel(), help)
}

// Listen for message on socket
func (b *Bot) Listen(ctx context.Context) {

	for {
		select {
		case ev := <-b.RTM.IncomingEvents:
			switch v := ev.Data.(type) {
			case *slack.MessageEvent:
				go b.process(NewMessage(v))

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", v.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
			}
		case <-ctx.Done():
			return
		}
	}
}

// Command adds a new command with custom handler
func (b *Bot) Command(cmd string, handler Handler) {
	b.Commands = append(b.Commands, NewCommand(b.CmdPrefix+cmd, "", handler))
}

// UnknownCommand will be called when the user calls a command that is unknown,
// but it will only work when the bot is in reply only mode
func (b *Bot) UnknownCommand(h Handler) {
	b.unknownCmdHandler = h
}

// Register registers a Command
func (b *Bot) Register(cmd CommandInterface) {
	b.Commands = append(b.Commands, cmd)
}
