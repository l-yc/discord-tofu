package docs

import (
	"github.com/bwmarrin/discordgo"
)

type MessageCreateHandler func(s *discordgo.Session, m *discordgo.MessageCreate)

type Command struct {
	Desc	string
	Fn		func(s *discordgo.Session, m *discordgo.MessageCreate)
}

type CommandDoc struct {
	Bind		string
	Command Command
}

var CommandDocList map[string][]CommandDoc
var HelpString string

func init() {
	CommandDocList = make(map[string][]CommandDoc)
}

func AddCommand(pkg string, bind string, cmd Command) {
	CommandDocList[pkg] = append(CommandDocList[pkg], CommandDoc{ Bind: bind, Command: cmd })
}

func CompileHelp() {
	HelpString = "__Help__\n"
	for pkg, cmdDocList := range CommandDocList {
		HelpString += "**" + pkg + ":**\n"
		for _, cmdDoc := range cmdDocList {
			HelpString += "* " + cmdDoc.Bind + ": " + cmdDoc.Command.Desc + "\n"
		}
	}
}

func GetHelp() string {
	return HelpString
}
