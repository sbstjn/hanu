# allot [![MIT License](https://img.shields.io/github/license/sbstjn/allot.svg?maxAge=3600)](https://github.com/sbstjn/allot/blob/master/LICENSE.md) [![GoDoc](https://godoc.org/github.com/sbstjn/allot?status.svg)](https://godoc.org/github.com/sbstjn/allot) [![Go Report Card](https://goreportcard.com/badge/github.com/sbstjn/allot)](https://goreportcard.com/report/github.com/sbstjn/allot) [![allot - Coverage Status](https://img.shields.io/coveralls/sbstjn/allot.svg)](https://coveralls.io/github/sbstjn/allot) [![Build Status](https://img.shields.io/circleci/project/sbstjn/allot.svg?maxAge=600)](https://circleci.com/gh/sbstjn/allot)

**allot** is a small `Golang` library to match and parse commands with pre-defined strings. For example use **allot** to define a list of commands your CLI application or Slackbot supports and check if incoming requests are matching your commands.

The **allot** library supports placeholders and regular expressions for parameter matching and parsing.

## Usage

```go
cmd := allot.NewCommand("revert <commits:integer> commits on <project:string> at (stage|prod)")
match, err := cmd.Match("revert 12 commits on example at prod")

if (err != nil)
  commits, _ = match.Integer("commits")
  project, _ = match.String("project")
  env, _ = match.Match(2)

  fmt.Printf("Revert \"%d\" on \"%s\" at \"%s\"", commits, project, env)
} else {
  fmt.Println("Request did not match command.")
}
```

## Examples

See the [hanu Slackbot](https://github.com/sbstjn/hanu) framework for a usecase for **allot**:

* [Host a Golang Slack bot on Heroku](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html)

## Credits
 * [Go coverage script from Mathias Lafeldt](https://mlafeldt.github.io/blog/test-coverage-in-go/)
