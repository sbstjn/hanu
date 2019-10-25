package hanu

import "testing"

func TestMessage(t *testing.T) {
	msg := Message{
		ID:      0,
		UserID:  "test",
		Type:    "message",
		Message: "text",
	}

	if msg.User() != "test" {
		t.Errorf("User() should be \"test\"")
	}

	if msg.Text() != "text" {
		t.Errorf("Text() should be \"test\", is \"%s\"", msg.Text())
	}

	if !msg.IsMessage() {
		t.Errorf("IsMessage() must be true")
	}

	if msg.IsDirectMessage() {
		t.Errorf("msg.IsDirectMessage() must be false")
	}

	if msg.IsMentionFor("") {
		t.Errorf("msg.IsMentionFor() must be false")
	}

	if msg.IsRelevantFor("user") {
		t.Errorf("msg.IsRelevantFor() must be true")
	}
}

func TestHelpMessage(t *testing.T) {
	msg := Message{}
	msg.SetText("help")

	if !msg.IsHelpRequest() {
		t.Errorf("msg.IsHelpRequest() must be true")
	}
}

func TestStripMention(t *testing.T) {
	msg := Message{}
	msg.SetText("<@test> help")

	msg.StripMention("test")

	if msg.Text() != "help" {
		t.Errorf("msg.Text must be 'help', is \"%s\"", msg.Text())
	}
}

func TestStripFormatting(t *testing.T) {
	var data = []struct {
		in  string
		out string
	}{
		{"how are you?", "how are you?"},
		{"<@test> how are you?", "how are you?"},
		{"<@test> Hi <http://example.com|example.com> test <https://lorem.ipsum|ipsum.com> fail", "Hi example.com test ipsum.com fail"},
		{"<@test> Hi <http://example.com|example.com> test <https://lorem.ipsum> fail", "Hi example.com test https://lorem.ipsum fail"},
		{"<@test> Hi <slackbot://example@domain.tld>", "Hi slackbot://example@domain.tld"},
		{"<@test> Hi <slackbot://example@domain.tld|label>", "Hi label"},
		{"<@U02UNKSJ1> is not changed", "<@U02UNKSJ1> is not changed"},
		{"<#C0C9MF8AK|channel> is not changed", "<#C0C9MF8AK|channel> is not changed"},
		{"<!subteam^S2N709QUT|@team> is not changed", "<!subteam^S2N709QUT|@team> is not changed"},
		{"asdas <http://google.com> asdasd", "asdas http://google.com asdasd"},
	}

	for _, set := range data {
		msg := Message{}
		msg.SetText(set.in)

		msg.StripMention("test")
		msg.StripLinkMarkup()

		if msg.Text() != set.out {
			t.Errorf("Failed to strip markup: \n Got: %s\n Expected: %s", msg.Text(), set.out)
		}
	}
}
