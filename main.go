package main

import (
	"os"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

const table = "┻━┻"

func main() {
	ListenForRage()
}

func ListenForRage() {
	bot := slackbot.New(os.Getenv("SLACK_API_TOKEN"))
	toMe := bot.Messages(slackbot.DirectMention, slackbot.DirectMessage)
	toMe.Hear(table).MessageHandler(SetTheTable)
	bot.Run()
}

func SetTheTable(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "┬─┬ ノ( ^_^ノ)", true)
}
