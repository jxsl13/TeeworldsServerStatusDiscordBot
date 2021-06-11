package main

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
)

type Config struct {
	DiscordToken     string
	DiscordChannelID discord.ChannelID
	Owner            string

	Servers         []string
	RefreshInterval time.Duration

	CustomFlags map[string]string
}

func (c *Config) Name() string {
	return path.Base(os.Args[0])
}

func discordChannelIDParser(out *discord.ChannelID) configo.ParserFunc {
	return func(value string) error {
		result, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = discord.ChannelID(result)
		return nil
	}
}

func (c *Config) Options() configo.Options {
	listDelimiter := ","
	pairDelimiter := ";"
	keyValueDelimiter := "->"
	return []configo.Option{
		{
			Key:           "DISCORD_TOKEN",
			Mandatory:     true,
			Description:   "Discord token retrieved from the discord developer portal.",
			ParseFunction: parsers.String(&c.DiscordToken),
		},
		{
			Key:           "DISCORD_CHANNEL_ID",
			Mandatory:     true,
			Description:   "Discord channel ID is a long integer",
			ParseFunction: discordChannelIDParser(&c.DiscordChannelID),
		},
		{
			Key:           "DICORD_OWNER",
			Description:   "username#1234 of the owner of this bot.",
			ParseFunction: parsers.String(&c.Owner),
		},
		{
			Key:           "TEEWORLDS_SERVERS",
			Mandatory:     true,
			Description:   "Comma separated list of server addresses e.g. 127.0.0.1:8303,127.0.0.1:8304",
			ParseFunction: parsers.List(&c.Servers, &listDelimiter),
		},
		{
			Key:           "REFRESH_INTERVAL",
			Mandatory:     true,
			Description:   "Interval until the server status message is refreshed again.",
			DefaultValue:  "10s",
			ParseFunction: parsers.Duration(&c.RefreshInterval),
		},
		{
			Key:           "CUSTOM_FLAGS",
			Description:   "Add custom flag mappings: 'default-><:default:853002988364103750>;XSC-><:default:853002988364103750>'",
			DefaultValue:  "default-><:default:853002988364103750>",
			ParseFunction: parsers.Map(&c.CustomFlags, &pairDelimiter, &keyValueDelimiter),
		},
	}
}
