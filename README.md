# hanu | [![MIT License](https://img.shields.io/github/license/sbstjn/hanu.svg?maxAge=3600)](https://github.com/sbstjn/hanu/blob/master/LICENSE.md) [![GoDoc](https://godoc.org/github.com/sbstjn/hanu?status.svg)](https://godoc.org/github.com/sbstjn/hanu) [![Go Report Card](https://goreportcard.com/badge/github.com/sbstjn/hanu)](https://goreportcard.com/report/github.com/sbstjn/hanu) [![Hanu - Coverage Status](https://img.shields.io/coveralls/sbstjn/hanu.svg)](https://coveralls.io/github/sbstjn/hanu) [![Build Status](https://img.shields.io/circleci/project/sbstjn/hanu.svg?maxAge=600)](https://circleci.com/gh/sbstjn/hanu)

The `Go` framework **hanu** is your best friend to create [Slack](https://slackhq.com) bots! **hanu** uses [platzhalter](https://github.com/sbstjn/platzhalter) for easy command parsing (e.g. `whisper <word>`) and runs fine as a [Heroku worker](https://devcenter.heroku.com/articles/background-jobs-queueing). All you need is a [Slack API token](https://api.slack.com/bot-users) and you can create your first bot within seconds! Just have a look at the [hanu-example](https://github.com/sbstjn/hanu-example) bot …

### Features

- Respond to **mentions**
- Respond to **direct messages**
- Auto-Generated command list for `help`
- Works fine as a **worker** on Heroku

## Usage

Use the following example code or the [hanu-example](https://github.com/sbstjn/hanu-example) bot to get started.

```go
package main

import (
	"log"
	"strings"

	"github.com/sbstjn/hanu"
	"github.com/sbstjn/hanu/conversation"
)

func main() {
	slack, err := hanu.New("SLACK_BOT_API_TOKEN")

	if err != nil {
		log.Fatal(err)
	}

	Version := "0.0.1"

	slack.Command("shout <word>", func(conv conversation.Interface) {
		conv.Reply(strings.ToUpper(conv.Param("word")))
	})

	slack.Command("whisper <word>", func(conv conversation.Interface) {
		conv.Reply(strings.ToLower(conv.Param("word")))
	})

	slack.Command("version", func(conv conversation.Interface) {
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

### Slack

Use direct messages for communication:

```
/msg @hanu version
```

Or use the bot in a public channel:

```
@hanu version
```


## Dependencies

 - [github.com/sbstjn/platzhalter](https://github.com/sbstjn/platzhalter) for parsing `cmd <key1> <key2>` strings
 - [golang.org/x/net/websocket](http://golang.org/x/net/websocket) for websocket communication with Slack

## Credits
 * [Host Go Slackbot on Heroku](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html)
 * [OpsDash article about Slack Bot](https://www.opsdash.com/blog/slack-bot-in-golang.html)
 * [Go coverage script from Mathias Lafeldt](https://mlafeldt.github.io/blog/test-coverage-in-go/)
