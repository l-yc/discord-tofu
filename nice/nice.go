package nice

import (
	"github.com/bwmarrin/discordgo"

	"errors"
	"strconv"

	"gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

var (
	PACKAGE = "Nice"
	DESC = "For the redditors."

	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	FnMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)

	db *gorm.DB
)

type User struct {
	gorm.Model
	ID				string
	Username  string
	NiceScore uint64	// TODO: are we going to create a db for each package?
}

func init() {
	WatchMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))
	FnMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))

	// Set up DB
	connectDB()

	WatchMap["nice"] = watchNice
	WatchMap["Nice"] = watchNice

	FnMap["register"] = registerUser
	FnMap["niceScore"] = niceScoreUser
	FnMap["niceScoreBoard"] = niceScoreBoard
}

func connectDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
    panic("failed to connect database")
  }
  db.AutoMigrate(&User{})
}

func watchNice(s *discordgo.Session, m *discordgo.MessageCreate) {
	var user User
	result := db.First(&user, m.Author.ID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.ChannelMessageSend(m.ChannelID, "who are you?")
	} else {
		user.NiceScore += 1
		db.Save(user)
		s.ChannelMessageSend(m.ChannelID, "nice")
	}
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
