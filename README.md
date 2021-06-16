# hanu - Go for Slack Bots!

[![Current Release](https://badgen.now.sh/github/release/sbstjn/hanu)](https://github.com/sbstjn/hanu/releases)
[![MIT License](https://badgen.now.sh/badge/License/MIT/blue)](https://github.com/sbstjn/hanu/blob/master/LICENSE.md)
[![Read Tutorial](https://badgen.now.sh/badge/Read/Tutorial/orange)](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html)
[![Code Example](https://badgen.now.sh/badge/Code/Example/cyan)](https://github.com/sbstjn/hanu-example)

The `Go` framework **hanu** is your best friend to create [Slack](https://slackhq.com) bots! **hanu** uses [allot](https://github.com/sbstjn/allot) for easy command and request parsing (e.g. `whisper <word>`) and runs fine as a [Heroku worker](https://devcenter.heroku.com/articles/background-jobs-queueing). All you need is a [Slack API token](https://api.slack.com/bot-users) and you can create your first bot within seconds! Just have a look at the [hanu-example](https://github.com/sbstjn/hanu-example) bot or [read my tutorial](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html) …

### Features

- Respond to **mentions**
- Respond to **direct messages**
- Auto-Generated command list for `help`
- Works fine as a **worker** on Heroku

## V1 Usage

Use the following example code or the [hanu-example](https://github.com/sbstjn/hanu-example) bot to get started.

```go
package main

import (
	"log"
	"strings"

	"github.com/sbstjn/hanu"
)

func main() {
	slack, err := hanu.New("SLACK_BOT_API_TOKEN")

	if err != nil {
		log.Fatal(err)
	}

	Version := "0.0.1"

	slack.Command("shout <word>", func(conv hanu.ConversationInterface) {
		str, _ := conv.String("word")
		conv.Reply(strings.ToUpper(str))
	})

	slack.Command("whisper <word>", func(conv hanu.ConversationInterface) {
		str, _ := conv.String("word")
		conv.Reply(strings.ToLower(str))
	})

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Thanks for asking! I'm running `%s`", Version)
	})

	slack.Listen()
}
```

The example code above connects to Slack using `SLACK_BOT_API_TOKEN` as the bot's token and can respond to direct messages and mentions for the commands `shout <word>` , `whisper <word>` and `version`.

You don't have to care about `help` requests, **hanu** has it built in and will respond with a list of all defined commands on direct messages like this:

```
/msg @hanu help
```

Of course this works fine with mentioning you bot's username as well:

```
@hanu help
```

Use direct messages for communication:

```
/msg @hanu version
```

Or use the bot in a public channel:

```
@hanu version
```

You can set the command prefix, if you like using those:

```go
bot.SetCommandPrefix("!")
bot.SetReplyOnly(false)
```

This will make it so you have to type:

```
!whisper I love turtles
```

For the command to be recognised.  Setting the bot to not reply only means it will listen to
all messages in an attempt to find a command (except help will only be printed when bot is mentioned).

Also, the `ConversationInterface` was changed to just `Convo` to save your wrists:

```go
	bot.Command("whisper <word>", func(conv hanu.Convo) {
		str, _ := conv.String("word")
		conv.Reply(strings.ToLower(str))
	})
```

The bot can also now talk arbitrarily and has a Channel object that is easy to
interface with since it's one function:

```go
bot.Say("UGHXISDF324", "I like %s", "turtles")

devops := bot.Channel("UGHXISDF324")
devops.Say("Host called %s is not responding to pings", "bobsburgers01")
```

You can print the help message whenever you want:

```go
bot.Say("UGHXISDF324", bot.BuildHelpText())
```

And there is an unknown command handler, but it only works when in reply only mode:

```go
bot.SetReplyOnly(true).UnknownCommand(func(c hanu.Convo) {
	c.Reply(slack.BuildHelpText())
})
```

Finally there is the ability to read messages that come into the channel in real time:

```go
devops := bot.Channel("UGHXISDF324")
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

for {
	select {
	case msg := <-devops.Messages():
		if msg.IsFrom("bob") {
			bot.Say("shutup <@bob>")
		}
	case <-ctx.Done():
		break
	}
}
```

## Dependencies

- [github.com/sbstjn/allot](https://github.com/sbstjn/allot) for parsing `cmd <param1:string> <param2:integer>` strings
- [golang.org/x/net/websocket](http://golang.org/x/net/websocket) for websocket communication with Slack
- [github.com/nlopes/slack](http://github.com/nlopes/slack) for real time communication with Slack

## Credits

- [Host Go Slackbot on Heroku](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html)
- [OpsDash article about Slack Bot](https://www.opsdash.com/blog/slack-bot-in-golang.html)
- [A Simple Slack Bot in Go - The Bot](ttps://dev.to/shindakun/a-simple-slack-bot-in-go---the-bot-4olg)