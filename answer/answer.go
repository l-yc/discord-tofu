package answer

import (
	"github.com/bwmarrin/discordgo"

	"github.com/l-yc/discord-tofu/config"
	"github.com/l-yc/discord-tofu/docs"

	"github.com/l-yc/discord-tofu/advice"
	"github.com/l-yc/discord-tofu/nice"
	"github.com/l-yc/discord-tofu/pics"

	"strings"
)

var (
	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	CmdMap	 map[string]docs.Command
)

func init() {
	WatchMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))
	CmdMap = make(map[string]docs.Command)

	for k, v := range advice.WatchMap {
		WatchMap[k] = v
	}
	for k, v := range advice.CmdMap {
		CmdMap[k] = v
		docs.AddCommand(advice.PACKAGE, k, v)
	}

	for k, v := range nice.WatchMap {
		WatchMap[k] = v
	}
	for k, v := range nice.CmdMap {
		CmdMap[k] = v
		docs.AddCommand(nice.PACKAGE, k, v)
	}

	for k, v := range pics.WatchMap {
		WatchMap[k] = v
	}
	for k, v := range pics.CmdMap {
		CmdMap[k] = v
		docs.AddCommand(pics.PACKAGE, k, v)
	}

	docs.CompileHelp()

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

	// Listen for commands that start with the prefix
	if len(m.Content) >= 2 && m.Content[0:2] == "::" {
		// Strip the prefix before handling
		m.Content = m.Content[2:]
		handleMessage(s, m)
	} else {
		// Otherwise, check if message exists in watchlist
		if watch, exists := WatchMap[m.Content]; exists {
			watch(s, m)
		}
	}
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")

	if cmd, exists := CmdMap[args[0]]; exists {
		cmd.Fn(s, m)
	} else {
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
			s.ChannelMessageSend(m.ChannelID, "こんにちは" + m.Author.Username + "-さん!")
			break
		case "whoami":
			s.ChannelMessageSend(m.ChannelID, "You are " + m.Author.ID)
			break
		case "help!":
			s.ChannelMessageSend(m.ChannelID, "がんばれ" + m.Author.Username + "！")
			break
		case "help":
			s.ChannelMessageSend(m.ChannelID, docs.GetHelp())
			break
		default:
			s.ChannelMessageSend(m.ChannelID, "なに？")
			break
		}
	}
}
