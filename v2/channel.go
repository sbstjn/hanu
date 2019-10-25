package hanu

// Channel is an object that allows a bot to say things without
// specifying the channel in every function call
type Channel struct {
	bot *Bot
	ID  string
}

// Say will cause the bot to say something in the channel
func (ch *Channel) Say(msg string, a ...interface{}) {
	ch.bot.Say(ch.ID, msg, a...)
}
