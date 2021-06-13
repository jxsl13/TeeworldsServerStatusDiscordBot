package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	abot "github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	bot "github.com/jxsl13/TeeworldsServerStatusDiscordBot/bot"
)

type Bot struct {
	Ctx *abot.Context
	Cfg *Config
}

// Stop stops the bot
func (b *Bot) Stop(event *gateway.MessageCreateEvent) (string, error) {
	if event.ChannelID != b.Cfg.DiscordChannelID {
		return "", nil
	}

	err := deleteMyMessages((*bot.Context)(b.Ctx), event.ChannelID)
	if err != nil {
		log.Println(err)
		return "", nil
	}

	err = b.Ctx.DeleteMessage(event.ChannelID, event.ID)
	if err != nil {
		log.Println(err)
		return "", nil
	}

	err = b.Ctx.CloseGracefully()
	if err != nil {
		log.Println(err)
		return "", nil
	}

	os.Exit(0)
	return "", nil
}

func deleteMyMessages(ctx *bot.Context, channelID discord.ChannelID) error {
	messages, err := ctx.Messages(channelID)
	if err != nil {
		log.Println(err)
		return nil
	}

	msgIDs := make([]discord.MessageID, 0, len(messages))
	for _, message := range messages {
		if message.Author.ID == ctx.Ready().User.ID {
			msgIDs = append(msgIDs, message.ID)
		}
	}

	err = ctx.DeleteMessages(channelID, msgIDs)
	if err != nil {

		return err
	}
	return nil
}

func NewRBACPrefix(config *Config) abot.Prefixer {
	return func(msg *gateway.MessageCreateEvent) (string, bool) {

		if strings.HasPrefix(msg.Content, "#") && config.Owner == msg.Author.Tag() {
			log.Printf("%s executed '%s'\n", msg.Author.Tag(), msg.Content)
			return "#", true
		}
		log.Printf("%s has no access to command '%s'\n", msg.Author.Tag(), msg.Content)
		return "", false
	}
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
