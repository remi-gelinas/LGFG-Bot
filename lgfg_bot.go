package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/remi-gelinas/lgfg-bot/internal/modules"
	"github.com/remi-gelinas/lgfg-bot/internal/router"
)

var Session, _ = discordgo.New()
var Router = router.New("!lgfg")

func init() {
	// Parse config from CLI flags
	flag.StringVar(&Session.Token, "token", "", "Discord auth token")

	// Add session event handlers
	Session.AddHandler(Router.OnMessageCreate)
	Session.AddHandler(modules.InLineAssignHandler)

	// Add router handlers
	Router.Route("help", "Ohai", nil, HelpHandler)

	Router.DefaultRoute(func(ds *discordgo.Session, msg *discordgo.Message, ctx *router.Context) {
		_, err := ds.ChannelMessageSend(msg.ChannelID, "You're going to have to speak my language so I can understand. Try "+Router.Prefix+" help to see a list of my commands.")

		if err != nil {
			fmt.Print(err)
		}
	})

	// Register command routes
	Router.Route("issue", "Ohai", []string{"string"}, modules.IssueHandler)
	Router.Route("feedback", "Nani", []string{"string"}, modules.FeedbackHandler)
}

func main() {
	var err error

	flag.Parse()

	// Verify token is not empty

	Session.Token = "Bot " + Session.Token

	// Open Discord socket connection
	err = Session.Open()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	Session.Close()
}

func HelpHandler(ds *discordgo.Session, msg *discordgo.Message, ctx *router.Context) {
	_, err := ds.ChannelMessageSend(msg.ChannelID, "```"+Router.Help+"```")

	if err != nil {
		// Failed to send message
	}
}
