package modules

import (
	"github.com/bwmarrin/discordgo"
)

func InLineAssignHandler(ds *discordgo.Session, gma *discordgo.GuildMemberAdd) {
	var err error

	guilds, err := ds.UserGuilds(1, "", "")

	if err != nil {
		// Failed to retrieve bot guilds
		return
	}

	var lgfgGuild *discordgo.UserGuild

	for _, g := range guilds {
		if g.Name == "Look Good Feel Good" {
			lgfgGuild = g
		}
	}

	if lgfgGuild == nil {
		// LGFG does not exist, panic
		return
	}

	lgfgRoles, err := ds.GuildRoles(lgfgGuild.ID)

	if err != nil {
		// Failed to get roles, panic
		return
	}

	var startRole *discordgo.Role

	for _, r := range lgfgRoles {
		if r.Name == "In Line" {
			startRole = r
		}
	}

	if startRole == nil {
		// Initial role not found, panic
		return
	}

	err = ds.GuildMemberRoleAdd(lgfgGuild.ID, gma.User.ID, startRole.ID)
}
