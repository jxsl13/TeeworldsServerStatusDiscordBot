package main

import (
	"fmt"
	"strings"

	"github.com/jxsl13/TeeworldsServerStatusDiscordBot/markdown"
	"github.com/jxsl13/twapi/browser"
)

func formatServerInfo(serverInfo browser.ServerInfo) string {

	var sb strings.Builder

	longestName := 0
	longestClan := 0
	for _, player := range serverInfo.Players {
		if len([]rune(player.Name)) > longestName {
			longestName = len([]rune(player.Name))
		}

		if len([]rune(player.Clan)) > longestName {
			longestName = len([]rune(player.Clan))
		}
	}

	header := fmt.Sprintf("%s (%d/%d)", markdown.WrapInInlineCodeBlock(serverInfo.Name), serverInfo.NumPlayers, serverInfo.MaxPlayers)

	sb.WriteString(header)
	sb.WriteString("\n")

	nameFmtStr := fmt.Sprintf("%%-%ds", longestName)
	clanFmtStr := fmt.Sprintf("%%-%ds", longestClan)

	for _, player := range serverInfo.Players {
		flag := markdown.Flag(player.Country)
		name := markdown.WrapInInlineCodeBlock(fmt.Sprintf(nameFmtStr, player.Name))
		clan := markdown.WrapInInlineCodeBlock(fmt.Sprintf(clanFmtStr, player.Clan))

		sb.WriteString(fmt.Sprintf("%s %s %s (%d)\n", flag, name, clan, player.Score))
	}

	return sb.String()
}
