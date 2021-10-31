package nice

import (
    "errors"
    "log"
    "os"
    "path/filepath"
    "strconv"

    "github.com/bwmarrin/discordgo"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"

    "github.com/l-yc/discord-tofu/docs"
    "github.com/l-yc/discord-tofu/config"
)

var (
	PACKAGE = "Nice"
	DESC = "For the redditors."

	WatchMap map[string]func (s *discordgo.Session, m *discordgo.MessageCreate)
	CmdMap map[string]docs.Command

	db *gorm.DB
)

type User struct {
	gorm.Model
	ID				string
	Username  string
	NiceScore uint64	// TODO: are we going to create a db for each package?
	Guilds		[]*Guild `gorm:"many2many:user_guilds"`
}

type Guild struct {
	gorm.Model
	ID string
	Users []*User `gorm:"many2many:user_guilds"`
}

func init() {
	WatchMap = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate))
	CmdMap = make(map[string]docs.Command)

	// Set up DB
	connectDB()

	WatchMap["nice"] = watchNice
	WatchMap["Nice"] = watchNice

	CmdMap["register"] = docs.Command{
		Desc: "Registers user in the database and server scoreboard.", Fn: registerUser }
	CmdMap["niceScore"] = docs.Command{
		Desc: "Displays nice score for the current user.", Fn: niceScoreUser }
	CmdMap["niceScoreBoard"] = docs.Command{
		Desc: "Displays nice scoreboard.", Fn: niceScoreBoard }
}

func connectDB() {
    dir := config.Cfg.DataDirectory
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        err = os.Mkdir(dir, 0755)
        if err != nil {
            log.Fatalf("Cannot create data directory: %v", err)
        }
        log.Printf("Created data directory at %s", dir)
    }

    file := filepath.Join(dir, "data.db")
    var err error
    db, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect database")
    }
    db.AutoMigrate(&User{})
    db.AutoMigrate(&Guild{})
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

		db.Model(&user).Association("Guilds").Append(&Guild{ ID: m.GuildID })
		s.ChannelMessageSend(m.ChannelID, "Updated your servers.")
	} else {
		db.Create(&User{
			ID: m.Author.ID,
			Username: m.Author.Username,
			NiceScore: 0,
			Guilds: []*Guild{&Guild{
				ID: m.GuildID,
			}},
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
	type Result struct {
		Username string
		NiceScore int
	}

	var results []Result

	db.Table("user_guilds").
		Select("users.username, users.nice_score").
		Where("guild_id = ?", m.GuildID).
		Joins("JOIN users on user_id = users.id").
		Order("nice_score desc").
		Limit(5).
		Scan(&results)

	msg := "Top " + strconv.Itoa(len(results)) + ":\n"
	for i, u := range results {
		msg += strconv.Itoa(i+1) + ". **" + u.Username + "**: "
		msg += strconv.Itoa(u.NiceScore) + " nices\n"
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}
