package modules

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/remi-gelinas/lgfg-bot/internal/router"
)

func IssueHandler(ds *discordgo.Session, msg *discordgo.Message, ctx *router.Context) {
	if len(ctx.Fields) == 1 {
		return
	}

	i := strings.Join(ctx.Fields[1:], " ")

	guilds, err := ds.UserGuilds(1, "", "")
	if err != nil {
		// Issue retrieving bot guilds
		return
	}

	var lgfgGuild *discordgo.UserGuild

	for _, g := range guilds {
		if g.Name == "Look Good Feel Good" {
			lgfgGuild = g
		}
	}

	if lgfgGuild == nil {
		// Something fucked up, LGFG doesn't exist
		return
	}

	lgfgChannels, err := ds.GuildChannels(lgfgGuild.ID)

	if err != nil {
		// Error retrieving LGFG channels
		return
	}

	var issuesChannel *discordgo.Channel

	for _, c := range lgfgChannels {
		if c.Name == "issues" {
			issuesChannel = c
		}
	}

	if issuesChannel == nil {
		// No issues channel present in LGFG
		return
	}

	ds.ChannelMessageSend(issuesChannel.ID, "**New issue from "+msg.Author.Username+"**: "+i)
}
