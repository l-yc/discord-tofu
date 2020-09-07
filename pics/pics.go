package pics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/l-yc/discord-tofu/config"
	"github.com/l-yc/discord-tofu/docs"

	"os"
	"log"
	"strings"
	"path"
)

var (
	PACKAGE = "Pics"
	DESC = "Ask tofu for some nice pics!"

	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	CmdMap map[string]docs.Command
)

func init() {
	WatchMap = make(map[string]func (s *discordgo.Session, m *discordgo.MessageCreate))
	CmdMap = make(map[string]docs.Command)

	CmdMap["react"] = docs.Command{
		Desc: "Tofu reacts.",
		Fn: func (s *discordgo.Session, m *discordgo.MessageCreate) {
			args := strings.Split(m.Content, " ")
			// default reaction
			var filename string = "GOMA.webp"
			if len(args) > 1 {
				// prevent traversal outside of the root pics directory
				filename = path.Clean(path.Join("/", args[1]))[1:]
			}

			file, err := os.Open(path.Join(config.Cfg.PicsDirectory, filename))
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, "I don't know how to do that ><")
			} else {
				_, err := s.ChannelFileSend(m.ChannelID, filename, file)
				if err != nil {
					log.Println(err)
				}
			}
		},
	}
}
