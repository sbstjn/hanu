package hanu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// Bot is the main object
type Bot struct {
	client            *socketmode.Client
	ID                string
	Commands          []CommandInterface
	ReplyOnly         bool
	CmdPrefix         string
	unknownCmdHandler Handler
	msgs              map[string]chan Message

	connectedWaiter chan bool
}

// New creates a new bot
func New(botToken, appToken string) (*Bot, error) {
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("bot token must have the prefix \"xoxb-\"")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, errors.New("app token must have the prefix \"xapp-\"")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	bot := &Bot{client: client, msgs: make(map[string]chan Message)}
	bot.connectedWaiter = make(chan bool)

	return bot, nil
}

// WaitForConnection will block until the bot is connected to the RTM
func (b *Bot) WaitForConnection() {
	if b.connectedWaiter == nil {
		return
	}
	<-b.connectedWaiter
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

func (b *Bot) notify(msg Message) {
	chnl := msg.Channel()
	ch, found := b.msgs[chnl]
	if !found {
		return
	}

	if cap(ch) == len(ch) {
		return
	}

	ch <- msg
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
	b.client.PostMessage(
		msg.Channel(),
		slack.MsgOptionText(msg.Text(), false),
	)
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

func (b *Bot) handleEvent(evt *socketmode.Event) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	fmt.Printf("Event received: %+v\n", eventsAPIEvent)

	b.client.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			_, _, err := b.client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				fmt.Printf("failed posting message: %v", err)
			}
			// router.HandleMention(ev)
		case *slackevents.MessageEvent:
			if os.Getenv("HANU_DEBUG") != "" {
				data, _ := json.MarshalIndent(evt, "", "  ")
				log.Println("NEW MSG ", string(data))
			}
			go b.process(NewMessage(ev))
			go b.notify(NewMessage(ev))
		case *slackevents.MemberJoinedChannelEvent:
			fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
		}
	default:
		b.client.Debugf("unsupported Events API event received")
	}
}

// Listen for message on socket
func (b *Bot) Listen(ctx context.Context) {
	for {
		select {
		case evt := <-b.client.Events:
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
				b.ID = evt.Request.ConnectionInfo.AppID
			case socketmode.EventTypeEventsAPI:
				b.handleEvent(&evt)
			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					// See https://api.slack.com/apis/connections/socket-implement#button

					b.client.Debugf("button clicked!")
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				b.client.Ack(*evt.Request, payload)
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				b.client.Debugf("Slash command received: %+v", cmd)

				payload := map[string]interface{}{
					"blocks": []slack.Block{
						slack.NewSectionBlock(
							&slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: "foo",
							},
							nil,
							slack.NewAccessory(
								slack.NewButtonBlockElement(
									"",
									"somevalue",
									&slack.TextBlockObject{
										Type: slack.PlainTextType,
										Text: "bar",
									},
								),
							),
						),
					},
				}

				b.client.Ack(*evt.Request, payload)
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
