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

## Usage

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

- [github.com/sbstjn/allot](https://github.com/sbstjn/allot) for parsing `cmd <param1:string> <param2:integer>` strings
- [golang.org/x/net/websocket](http://golang.org/x/net/websocket) for websocket communication with Slack

## Credits

- [Host Go Slackbot on Heroku](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html)
- [OpsDash article about Slack Bot](https://www.opsdash.com/blog/slack-bot-in-golang.html)
- [Go coverage script from Mathias Lafeldt](https://mlafeldt.github.io/blog/test-coverage-in-go/)
