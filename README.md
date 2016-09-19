# hanu | [![MIT License](https://img.shields.io/github/license/sbstjn/hanu.svg?maxAge=3600)](https://github.com/sbstjn/hanu/blob/master/LICENSE.md) [![GoDoc](https://godoc.org/github.com/sbstjn/hanu?status.svg)](https://godoc.org/github.com/sbstjn/hanu) [![Go Report Card](https://goreportcard.com/badge/github.com/sbstjn/hanu)](https://goreportcard.com/report/github.com/sbstjn/hanu) [![Coverage Status](https://coveralls.io/repos/github/sbstjn/hanu/badge.svg)](https://coveralls.io/github/sbstjn/hanu) [![Build Status](https://travis-ci.org/sbstjn/hanu.svg?branch=master)](https://travis-ci.org/sbstjn/hanu)

The `Go` framework **hanu** supports you when creating [Slack](https://slackhq.com) bots.

## Dependencies

 - [github.com/sbstjn/platzhalter](https://github.com/sbstjn/platzhalter) for parsing `cmd <key1> <key2>` strings
 - [golang.org/x/net/websocket](http://golang.org/x/net/websocket) for websocket communication with Slack

## Configuration

You need to create an [API token in Slack](https://api.slack.com/bot-users) for your *hanu* bot first. See the example project [hanu-example](https://github.com/sbstjn/hanu-example) for an example usage of environment variable and configuration YAML file.

## Usage

Use the following example code or use the [hanu-example](https://github.com/sbstjn/hanu-example) project to get started.

```go
package main

import (
	"log"
	"strings"

	"github.com/sbstjn/hanu"
)

func main() {
	bot, err := hanu.New("SLACK_BOT_API_TOKEN")

	if err != nil {
		log.Fatal(err)
	}

	Version := "0.0.1"

	bot.Command("shout <word>", func(conv *hanu.Conversation) {
		conv.Reply(strings.ToUpper(conv.Param("word")))
	})

	bot.Command("whisper <word>", func(conv *hanu.Conversation) {
		conv.Reply(strings.ToLower(conv.Param("word")))
	})

	bot.Command("version", func(conv *hanu.Conversation) {
		conv.Reply("Thanks for asking! I'm running `%s`", Version)
	})

	bot.Listen()
}
```

The example code connects to Slack using `SLACK_BOT_API_TOKEN` as the bot's token and if everythings works fine, your bot responds to direct messages and mentions for the command `shout <word>` , `whisper <word>` and `version`.

Use direct messages for communication:

```
/msg @hanu-example version
```

Or use the bot in a public channel:

```
@hanu-example version
```

## Credits
 * [OpsDash article about Slack Bot](https://www.opsdash.com/blog/slack-bot-in-golang.html)
 * [Go coverage script from Mathias Lafeldt](https://mlafeldt.github.io/blog/test-coverage-in-go/)
