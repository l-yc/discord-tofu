package answer

import (
	"github.com/bwmarrin/discordgo"

	"github.com/l-yc/discord-tofu/config"
	"github.com/l-yc/discord-tofu/advice"
	"github.com/l-yc/discord-tofu/nice"

	"strings"
)

var (
	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	FnMap		 map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	Help		 map[string][]string // TODO please use a struct
)

func init() {
	WatchMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))
	FnMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))
	Help = make(map[string][]string)

	for k, v := range advice.WatchMap {
		WatchMap[k] = v
	}
	for k, v := range advice.FnMap {
		FnMap[k] = v
		Help[advice.PACKAGE] = append(Help[advice.PACKAGE], k)
	}

	for k, v := range nice.WatchMap {
		WatchMap[k] = v
	}
	for k, v := range nice.FnMap {
		FnMap[k] = v
		Help[nice.PACKAGE] = append(Help[nice.PACKAGE], k)
	}

	WatchMap["<3"] = func (s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == config.Cfg.Owner {
			s.ChannelMessageSend(m.ChannelID, "I love you too <3")
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if message exists in watchlist
	if watch, exists := WatchMap[m.Content]; exists {
		watch(s, m)
		return
	}

	// Otherwise, ignore messages that doesn't start with the prefix
	if len(m.Content) < 2 || m.Content[0:2] != "::" {
		return
	}
	// Strip the prefix
	m.Content = m.Content[2:]
	// Handle the message!
	handleMessage(s, m)
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")

	if fn, exists := FnMap[args[0]]; exists {
		fn(s, m)
		return
	}

	switch args[0] {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		break
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
		break
	case "poke": // variety of health checks
		s.ChannelMessageSend(m.ChannelID, "poke")
		break
	case "hello": // basic functionality
		s.ChannelMessageSend(m.ChannelID, "こんにちは, " + m.Author.Username + "-さん!")
		break
	case "whoami":
		s.ChannelMessageSend(m.ChannelID, "You are " + m.Author.ID)
		break
	case "help":
		displayHelp(s, m)
	}
}

func displayHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "__Help__\n"
	for k, v := range Help {
		msg += "**" + k + ":**\n"
		for _, s := range v {
			msg += "* " + s + "\n"
		}
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}
