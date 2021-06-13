package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/jxsl13/TeeworldsServerStatusDiscordBot/bot"
	"github.com/jxsl13/TeeworldsServerStatusDiscordBot/markdown"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/twapi/browser"
)

var (
	cfg             = &Config{}
	scheduler       = gocron.NewScheduler(time.UTC)
	addressMsgIDMap = make(map[string]discord.MessageID)
)

func init() {
	useEnvFile := flag.String("env-file", "", "--env-file .env")
	flag.Parse()

	var env map[string]string
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
}

func main() {

	bot.RunWithShutdownCallback(cfg.DiscordToken, &Bot{Cfg: cfg},
		func(ctx *bot.Context) error {
			// prefix handling
			ctx.HasPrefix = NewRBACPrefix(cfg)

			// init messages for all servers
			for _, address := range cfg.Servers {
				// imporant line, don't touch
				address := address

				// create message for every address
				msg, err := ctx.SendMessage(cfg.DiscordChannelID, address, nil)
				if err != nil {
					log.Printf("skipping address %s due to error: %v\n", address, err)
					continue
				}

				// save initialized message ids
				addressMsgIDMap[address] = msg.ID

				ip, port, err := toIPAndPort(address)
				if err != nil {
					return err
				}

				// create a scheduler for every message
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
						cfg.DiscordChannelID, msg.ID, ip, port)
			}

			// start scheduler
			scheduler.
				StartAt(time.Now().Add(time.Second)).
				StartAsync()
			return nil
		},
		func(ctx *bot.Context) {
			// shutdown callback
			deleteMyMessages(ctx, cfg.DiscordChannelID)
		},
	)
}
