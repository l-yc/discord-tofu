package advice

import (
	"github.com/bwmarrin/discordgo"

	"math/rand"
	"time"
)

var (
	PACKAGE = "Advice"
	DESC = "Get advice from tofu!"

	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	FnMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)

	AdviceAnswers = []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
)

func init() {
	WatchMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))
	FnMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))

	// Seeding with the same value results in the same random sequence each run.
	// For different numbers, seed with a different value, such as
	// time.Now().UnixNano(), which yields a constantly-changing number.
	rand.Seed(time.Now().UnixNano())

	FnMap["advice"] = func (s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, AdviceAnswers[rand.Intn(len(AdviceAnswers))])
	}
}
