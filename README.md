# hanu | ![MIT License](https://img.shields.io/github/license/sbstjn/hanu.svg?maxAge=3600) [![GoDoc](https://godoc.org/github.com/sbstjn/hanu?status.svg)](https://godoc.org/github.com/sbstjn/hanu) [![Go Report Card](https://goreportcard.com/badge/github.com/sbstjn/hanu)](https://goreportcard.com/report/github.com/sbstjn/hanu) [![Coverage Status](https://coveralls.io/repos/github/sbstjn/hanu/badge.svg)](https://coveralls.io/github/sbstjn/hanu) [![Build Status](https://travis-ci.org/sbstjn/hanu.svg?branch=master)](https://travis-ci.org/sbstjn/hanu)

The `Go` framework **hanu** supports you when creating [Slack](https://slackhq.com) bots.

## Dependencies

 - [github.com/sbstjn/platzhalter](https://github.com/sbstjn/platzhalter) for parsing `cmd <key1> <key2>` strings
 - [golang.org/x/net/websocket](http://golang.org/x/net/websocket) for websocket communication with Slack

## Configuration

You need to create an [API token in Slack](https://api.slack.com/bot-users) for your *hanu* instance.

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

## Credits

The **hanu** framework uses the [coverage script from Mathias Lafeldt](https://mlafeldt.github.io/blog/test-coverage-in-go/)
