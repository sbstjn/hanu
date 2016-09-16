# hanu [![Build Status](https://travis-ci.org/sbstjn/hanu.svg?branch=master)](https://travis-ci.org/sbstjn/hanu) [![Coverage Status](https://coveralls.io/repos/github/sbstjn/hanu/badge.svg)](https://coveralls.io/github/sbstjn/hanu)

Create a `golang` Slack Bot with **hanu** and a few lines of code. Define commands and interact with user on Slack.

## Usage

```go
package main

import (
	"log"
	"strings"

	"github.com/sbstjn/hanu"
)

func main() {
	bot, err := hanu.New("SECRET")

	if err != nil {
		log.Fatal(err)
	}

	bot.Register("shout <word>", func(conv *hanu.Conversation) {
		conv.Reply(strings.ToUpper(conv.Param("word")))
	})

	bot.Register("whisper <word>", func(conv *hanu.Conversation) {
		conv.Reply(strings.ToLower(conv.Param("word")))
	})

	bot.Register("reverse <word>", func(conv *hanu.Conversation) {
		conv.Reply("Not available :(")
	})

	bot.Listen()
}
```
