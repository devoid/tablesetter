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

const zenCap = 8

type forceTracker struct {
	mu      sync.Mutex
	angry   map[string]time.Time
	zenmode map[string]int
	zenCap  int // after so many angry replies, revert to zen
}

func (ft *forceTracker) HandleCooldown(tick time.Duration, cooldown time.Duration) {
	ticker := time.NewTicker(tick)
	for {
		<-ticker.C
		now := time.Now()
		ft.mu.Lock()
		for u, t := range ft.angry {
			if t.Add(cooldown).Before(now) {
				delete(ft.angry, u)

				// clear zenmode flag
				if _, ok := ft.zenmode[u]; ok {
					delete(ft.zenmode, u)
				}
			}
		}
		ft.mu.Unlock()

	}
}

func (ft *forceTracker) IsUserAnnoying(user string) bool {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	ft.zenmode[user]++

	if _, ok := ft.angry[user]; ok && ft.zenmode[user] < ft.zenCap {
		log.Printf("User %s is annoying (cap %d)", user, ft.zenmode[user])
		ft.angry[user] = time.Now()
		return true
	}
	ft.angry[user] = time.Now()
	return false
}

func newForceTracker(zenCap int) *forceTracker {
	return &forceTracker{
		mu:      sync.Mutex{},
		angry:   make(map[string]time.Time),
		zenmode: make(map[string]int),
		zenCap:  zenCap,
	}
}

var tracker *forceTracker

func main() {
	tracker = newForceTracker(zenCap)
	go tracker.HandleCooldown(time.Second, time.Minute*5)

	bot := slackbot.New(os.Getenv("SLACK_API_TOKEN"))
	bot.MessageHandler(setTheTable)
	bot.Run()
}

func setTheTable(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	if strings.Contains(evt.Text, tables[0]) || strings.Contains(evt.Text, tables[1]) {
		resp := getReplyString(ctx, bot, evt)
		bot.Reply(evt, resp, true)
	}
}

func getReplyString(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) string {
	if tracker.IsUserAnnoying(evt.User) {
		return randString(angrySetter)
	}
	return randString(zenSetter)
}

func randString(set []string) string {
	r := rand.Intn(len(set))
	return set[r]
}
