package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

var tables = []string{
	"┻━┻",
	"┻┻",
}

var zenSetter = []string{
	"┳━┳ ノ( ^_^ノ)",
	"┳━┳ ノ(._.ノ)",
	"┳━┳ ¯\\_(ツ)",
	"┬──┬ ノ( ゜-゜ノ)",
	"(ヘ･_･)ヘ┳━┳",
}

var angrySetter = []string{
	"┬─┬ ╯(°□° ╯)",
	"┬─┬ ノ(ಥ益ಥノ",
	"┬─┬ ヽ(`Д´)ﾉ",
	"┬─┬ ヽ༼ຈل͜ຈ༽ﾉ",
	"┬─┬ ノ(TДTノ)",
}

type forceTracker struct {
	mu    sync.Mutex
	users map[string]time.Time
}

func (ft *forceTracker) HandleCooldown(tick time.Duration, cooldown time.Duration) {
	ticker := time.NewTicker(tick)
	for {
		<-ticker.C
		now := time.Now()
		ft.mu.Lock()
		for u, t := range ft.users {
			if t.Add(cooldown).Before(now) {
				delete(ft.users, u)
			}
		}
		ft.mu.Unlock()

	}
}

func (ft *forceTracker) IsUserAnnoying(user string) bool {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	if _, ok := ft.users[user]; ok {
		log.Printf("User %s is annoying", user)
		ft.users[user] = time.Now()
		return true
	}
	ft.users[user] = time.Now()
	return false
}

func newForceTracker() *forceTracker {
	return &forceTracker{mu: sync.Mutex{}, users: make(map[string]time.Time)}
}

var tracker *forceTracker

func main() {
	tracker = newForceTracker()
	go tracker.HandleCooldown(time.Second, time.Minute*5)

	bot := slackbot.New(os.Getenv("SLACK_API_TOKEN"))
	bot.MessageHandler(SetTheTable)
	bot.Run()
}

func SetTheTable(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	if strings.ContainsAny(evt.Text, tables) {
		resp := getReplyString(ctx, bot, evt)
		bot.Reply(evt, resp, true)
	}
}

func getReplyString(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) string {
	if tracker.IsUserAnnoying(evt.User) {
		return randString(angrySetter)
	} else {
		return randString(zenSetter)
	}
}

func randString(set []string) string {
	r := rand.Intn(len(set))
	return set[r]
}
