package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/jxsl13/TeeworldsServerStatusDiscordBot/markdown"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/twapi/browser"
)

var (
	cfg             = &Config{}
	scheduler       = gocron.NewScheduler(time.UTC)
	addressMsgIDMap map[string]discord.MessageID
)

func init() {
	useEnvFile := flag.String("env-file", "", "--env-file .env")
	flag.Parse()

	env := make(map[string]string)
	var err error

	if *useEnvFile != "" {
		env, err = godotenv.Read(*useEnvFile)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		env = configo.GetEnv()
	}

	if err := configo.Parse(cfg, env); err != nil {
		log.Fatalln(err)
	}

	scheduler.TagsUnique()
	addressMsgIDMap = make(map[string]discord.MessageID)
}

func main() {

	bot.Run(cfg.DiscordToken, &Bot{Cfg: cfg, AddressMsgIDMap: addressMsgIDMap},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = NewRBACPrefix(cfg)

			for _, address := range cfg.Servers {
				// imporant line, don't touch
				address := address

				// create message for every address
				msg, err := ctx.SendMessage(cfg.DiscordChannelID, address, nil)
				if err != nil {
					log.Printf("skipping address %s due to error: %v\n", address, err)
					continue
				}

				channelID := cfg.DiscordChannelID
				msgID := msg.ID

				addressMsgIDMap[address] = msgID

				ip, port, err := toIPAndPort(address)
				if err != nil {
					return err
				}

				scheduler.
					Every(cfg.RefreshInterval).
					Tag(address).
					Do(
						func(channelID discord.ChannelID, msgID discord.MessageID, ip string, port int) {

							updateMessage := ""

							serverInfo, err := browser.GetServerInfoWithTimeout(ip, port, cfg.RefreshInterval)
							if err != nil {
								updateMessage = fmt.Sprintf("%s: %s", markdown.WrapInFat("[DOWN]"), address)
								log.Println(err)
							} else {
								updateMessage = formatServerInfo(serverInfo)
							}

							msg, err = ctx.EditMessage(channelID, msgID, updateMessage, nil, true)
							if err != nil {
								log.Printf("could not edit message: %v\v", err)
								return
							}

						},
						channelID, msgID, ip, port)
			}

			scheduler.
				StartAt(time.Now().Add(time.Second)).
				StartAsync()
			return nil
		},
	)

}

var splitRegex = regexp.MustCompile(`(.+):(\d+)`)

func toIPAndPort(address string) (string, int, error) {
	matches := splitRegex.FindStringSubmatch(address)
	if len(matches) == 0 {
		return "", 0, fmt.Errorf("invalid address: %s", address)
	}
	ip := matches[1]
	port, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, err
	}

	return ip, port, nil
}

type Bot struct {
	Ctx             *bot.Context
	Cfg             *Config
	AddressMsgIDMap map[string]discord.MessageID
}

func (b *Bot) Stop(event *gateway.MessageCreateEvent) (string, error) {
	for _, msgID := range b.AddressMsgIDMap {
		b.Ctx.DeleteMessage(b.Cfg.DiscordChannelID, msgID)
	}
	b.Ctx.DeleteMessage(event.ChannelID, event.Message.ID)

	b.Ctx.CloseGracefully()
	os.Exit(0)
	return "", nil
}

func NewRBACPrefix(config *Config) bot.Prefixer {
	return func(msg *gateway.MessageCreateEvent) (string, bool) {

		if strings.HasPrefix(msg.Content, "#") && config.Owner == msg.Author.Tag() {
			return "#", true
		}
		return "", false
	}
}
