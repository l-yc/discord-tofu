package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/l-yc/discord-tofu/config"
	"github.com/l-yc/discord-tofu/answer"
)

type Flags struct {
	ConfigFile	string
}

var (
	flags	Flags
)

func init() {
	logFile := strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_") + ".log"
	f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	flag.StringVar(&flags.ConfigFile, "c", "config.toml", "Config File")
	flag.Parse()
}

func main() {
	config.ReadConfig(flags.ConfigFile)

	discord, err := discordgo.New("Bot " + config.Cfg.Token)
	if err != nil {
		log.Println("Error creating discord session:", err)
		return
	}
	// Invite the bot
	fmt.Println("Invite the bot to your server by visiting the following link:")
	fmt.Printf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot", config.Cfg.ClientID)
	fmt.Println()

	// Register the MessageCreate func as a callback for MessageCreate events.
	discord.AddHandler(answer.MessageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
