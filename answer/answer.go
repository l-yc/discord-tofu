package answer

import (
	"github.com/bwmarrin/discordgo"

	"errors"
	"strings"
	"strconv"
	"math/rand"
	"time"

	"gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

type User struct {
	gorm.Model
	ID				string
	Username  string
	NiceScore uint64
}

var (
	db *gorm.DB
	adviceAnswers = []string{
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
	// Seeding with the same value results in the same random sequence each run.
	// For different numbers, seed with a different value, such as
	// time.Now().UnixNano(), which yields a constantly-changing number.
	rand.Seed(time.Now().UnixNano())

	// for persisting data
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
    panic("failed to connect database")
  }

	// Migrate the schema
  db.AutoMigrate(&User{})
}

func watchFor(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	ret := false

	if strings.ToLower(m.Content) == "nice" {
		var user User
		result := db.First(&user, m.Author.ID)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.ChannelMessageSend(m.ChannelID, "who are you?")
		} else {
			user.NiceScore += 1
			db.Save(user)
			s.ChannelMessageSend(m.ChannelID, "nice")
			ret = true
		}
	}

	return ret
}

func registerUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	var user User
	result := db.First(&user, m.Author.ID)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if user.Username != m.Author.Username {
			user.Username = m.Author.Username
			db.Save(&user)
			s.ChannelMessageSend(m.ChannelID, "Got it. Nice name, " + m.Author.Username + "!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "I already know you!")
		}
	} else {
		db.Create(&User{
			ID: m.Author.ID,
			Username: m.Author.Username,
			NiceScore: 0,
		})
		s.ChannelMessageSend(m.ChannelID, "Alright, registered " + m.Author.Username + "!")
	}
}

func niceScoreUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	var user User
	result := db.First(&user, m.Author.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.ChannelMessageSend(m.ChannelID, "who are you?")
		return
	}

	msg := "Score for " + m.Author.Username + " = " + strconv.FormatUint(user.NiceScore, 10)
	s.ChannelMessageSend(m.ChannelID, msg)
}

func niceScoreBoard(s *discordgo.Session, m *discordgo.MessageCreate) {
	var users []User
	db.Order("nice_score desc").Find(&users).Limit(5)

	msg := "Top " + strconv.Itoa(len(users)) + ":\n"
	for i, u := range users {
		msg += strconv.Itoa(i+1) + ". **" + u.Username + "**: "
		msg += strconv.FormatUint(u.NiceScore, 10) + " nices\n"
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	switch args[0] {
	// If the message is "ping" reply with "Pong!"
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		break
	// If the message is "pong" reply with "Ping!"
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
		break
	case "advice":
		s.ChannelMessageSend(m.ChannelID, adviceAnswers[rand.Intn(len(adviceAnswers))])
		break
	case "hello":
		s.ChannelMessageSend(m.ChannelID, "こんにちは, " + m.Author.Username + "-さん!")
		break
	case "whoami":
		s.ChannelMessageSend(m.ChannelID, "You are " + m.Author.ID)
		break
	case "register":
		registerUser(s, m)
	case "niceScore":
		niceScoreUser(s, m)
		break
	case "niceScoreBoard":
		niceScoreBoard(s, m)
		break
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
	if watchFor(s, m) {
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
