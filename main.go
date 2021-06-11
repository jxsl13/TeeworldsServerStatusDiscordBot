package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/twapi/browser"
)

var (
	cfg           = &Config{}
	bulletinBoard = (*BulletinBoard)(nil)
)

func init() {
	if err := configo.Parse(cfg, configo.GetEnv()); err != nil {
		log.Fatalln(err)
	}

	bulletinBoard = NewBulletinBoard(uint64(cfg.DiscordChannelID))
}

func main() {

	bot.Run(cfg.DiscordToken, &Bot{Cfg: cfg},
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

				bulletinBoard.Register(address, msg.ID, cfg.RefreshInterval, func() error {

					address := address

					channelID := cfg.DiscordChannelID
					msgID := msg.ID

					ip, port, err := toIPAndPort(address)
					if err != nil {
						return err
					}
					serverInfo, err := browser.GetServerInfoWithTimeout(ip, port, cfg.RefreshInterval)
					if err != nil {
						return err
					}

					msg, err = ctx.EditMessage(channelID, msgID, formatServerInfo(serverInfo), nil, true)
					if err != nil {
						return err
					}

					bulletinBoard.UpdateMsgID(address, msg.ID)

					return nil
				})
			}

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
	Ctx *bot.Context
	Cfg *Config
	BB  *BulletinBoard
}

func (b *Bot) Stop(event *gateway.MessageCreateEvent) (string, error) {
	for _, msgID := range bulletinBoard.AddressMessageMap {
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
