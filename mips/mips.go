package mips

import (
	"github.com/bwmarrin/discordgo"
	"github.com/l-yc/discord-tofu/docs"

	"fmt"
	"strings"
	"strconv"
)

var (
	PACKAGE = "Mips"
	DESC = "Tofu takes CS2100!"

	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	CmdMap map[string]docs.Command
)

func init() {
	WatchMap = make(map[string]func (s *discordgo.Session, m *discordgo.MessageCreate))
	CmdMap = make(map[string]docs.Command)

	CmdMap["mips"] = docs.Command{
		Desc: "Tofu reads mips code. [enc|dec|doc]",
		Fn: func (s *discordgo.Session, m *discordgo.MessageCreate) {
			args := strings.Split(m.Content, " ")
			if len(args) < 3 {
				s.ChannelMessageSend(m.ChannelID, "Bzzzt.")
			} else {
				switch args[1] {
				case "enc":
					s.ChannelMessageSend(m.ChannelID, "my brain is too smol :(")
				case "dec":
					success := true
					code := ""

					hex := strings.Join(args[2:], " ") // join back
					hex = strings.TrimPrefix(hex, "0x")
					dec, err := strconv.ParseInt(hex, 16, 32)
					fmt.Println(err)
					success = success && (err == nil)

					bin := strconv.FormatInt(dec, 2)
					for len(bin) < 32 {
						bin = "0" + bin
					}

					fmt.Println(hex, dec, bin)
					op := bin[0:6]
					code = fmt.Sprintf("op code = %s, bin = %s", op, bin)

					if success {
						s.ChannelMessageSend(m.ChannelID, code)
					} else {
						s.ChannelMessageSend(m.ChannelID, "bad")
					}
				case "doc":
					s.ChannelMessageSend(m.ChannelID, "go google for now")
				default:
					s.ChannelMessageSend(m.ChannelID, "a")
				}
			}
		},
	}
}
