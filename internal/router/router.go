package router

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Route struct {
	Pattern     string
	Description string
	Args        []string
	Run         HandlerFunc
}

type Context struct {
	Fields          []string
	Content         string
	GuildID         string
	IsPrivate       bool
	IsDirected      bool
	HasMention      bool
	HasMentionFirst bool
	HasPrefix       bool
}

type HandlerFunc func(*discordgo.Session, *discordgo.Message, *Context)

type Router struct {
	Routes  []*Route
	Default *Route
	Help    string
	Prefix  string
}

func New(pre string) *Router {
	r := &Router{}
	r.Help = ""
	r.Prefix = pre

	return r
}

func (rr *Router) DefaultRoute(cb HandlerFunc) (*Route, error) {
	r := Route{}
	r.Run = cb
	rr.Default = &r

	return &r, nil
}

func (rr *Router) Route(pattern string, desc string, args []string, cb HandlerFunc) (*Route, error) {
	r := Route{}
	r.Pattern = pattern
	r.Description = desc
	r.Run = cb
	rr.Routes = append(rr.Routes, &r)

	// Build help string for the given Route
	argStr := ""

	for _, a := range args {
		argStr = argStr + "<" + a + "> "
	}

	argStr = strings.TrimSuffix(argStr, " ")

	if argStr != "" {
		argStr = " " + argStr
	}

	rr.Help = rr.Help + pattern + argStr + ": " + desc + "\n"

	return &r, nil
}

func (rr *Router) Match(msg string) (*Route, []string) {
	fields := strings.Fields(msg)

	if len(fields) == 0 {
		return nil, nil
	}

	for fk, fv := range fields {
		for _, rv := range rr.Routes {
			if rv.Pattern == fv {
				return rv, fields[fk:]
			}
		}
	}

	return nil, nil
}

func (rr *Router) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {
	var err error

	if mc.Author.ID == ds.State.User.ID {
		return
	}

	ctx := &Context{
		Content: strings.TrimSpace(mc.Content),
	}

	var c *discordgo.Channel
	c, err = ds.State.Channel(mc.ChannelID)

	if err != nil {
		c, err = ds.Channel(mc.ChannelID)
		if err != nil {
			// Unable to fetch channel for message
		} else {
			err = ds.State.ChannelAdd(c)
			if err != nil {
				// Error updating channel state
			}

			ctx.GuildID = c.GuildID
			if c.Type == discordgo.ChannelTypeDM {
				ctx.IsPrivate = true
				ctx.IsDirected = true
			}
		}
	}

	if !ctx.IsDirected {
		for _, v := range mc.Mentions {
			if v.ID == ds.State.User.ID {

				ctx.IsDirected, ctx.HasMention = true, true

				reg := regexp.MustCompile(fmt.Sprintf("<@!?(%s)>", ds.State.User.ID))

				if reg.FindStringIndex(ctx.Content)[0] == 0 {
					ctx.HasMentionFirst = true
				}

				ctx.Content = reg.ReplaceAllString(ctx.Content, "")

				break
			}
		}
	}

	if !ctx.IsDirected && len(rr.Prefix) > 0 {
		if strings.HasPrefix(ctx.Content, rr.Prefix) {
			ctx.IsDirected, ctx.HasPrefix, ctx.HasMentionFirst = true, true, true
			ctx.Content = strings.TrimPrefix(ctx.Content, rr.Prefix)
		}
	}

	// Only responds to prefixed commands
	if !ctx.HasPrefix {
		return
	}

	r, fl := rr.Match(ctx.Content)
	if r != nil {
		ctx.Fields = fl
		r.Run(ds, mc.Message, ctx)

		return
	}

	// Run default route, if existed
	if rr.Default != nil {
		rr.Default.Run(ds, mc.Message, ctx)
	}
}
